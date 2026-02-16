package timer

import (
	"context"
	"fmt"
	"time"

	"terminal-brain/internal/core"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

type Module struct {
	app       *core.AppState
	logs      []core.TimeLog
	activeID  int64
	activeLbl string
	label     textinput.Model
	starting  bool
	errMsg    string
	status    string
	cursor    int
	confirm   bool
}

type tickMsg time.Time

func New() *Module {
	label := textinput.New()
	label.Placeholder = "Task label"
	return &Module{label: label}
}

func (m *Module) Name() string {
	return "Timer"
}

func (m *Module) Help() string {
	return "Timer: s start | t stop | n/p move | d delete | D clear"
}

func (m *Module) InputActive() bool {
	return m.starting
}

func (m *Module) Init(app *core.AppState) error {
	m.app = app
	if err := m.reload(); err != nil {
		return err
	}
	return nil
}

func (m *Module) Update(msg tea.Msg) (core.Module, tea.Cmd) {
	if m.starting {
		return m.updateStart(msg)
	}

	switch msg := msg.(type) {
	case tickMsg:
		if m.activeID != 0 {
			return m, tickCmd()
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			if m.activeID == 0 {
				m.starting = true
				m.label.SetValue("")
				m.label.Focus()
				return m, nil
			}
		case "t":
			if m.activeID != 0 {
				m.stopTimer()
				return m, nil
			}
		case "n":
			if m.cursor < len(m.logs)-1 {
				m.cursor++
			}
			m.confirm = false
			return m, nil
		case "p":
			if m.cursor > 0 {
				m.cursor--
			}
			m.confirm = false
			return m, nil
		case "d":
			m.deleteSelected()
			return m, nil
		case "D":
			if m.confirm {
				m.clearAll()
				return m, nil
			}
			m.confirm = true
			m.status = "Press D again to clear all timers"
			return m, nil
		case "esc":
			m.confirm = false
			return m, nil
		case "r":
			_ = m.reload()
			return m, nil
		}
	}

	if m.activeID != 0 {
		return m, tickCmd()
	}
	return m, nil
}

func (m *Module) updateStart(msg tea.Msg) (core.Module, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.starting = false
			return m, nil
		case "enter":
			label := m.label.Value()
			if label == "" {
				label = "Untitled"
			}
			log, err := m.app.Stores.Time.Start(context.Background(), label, time.Now().UTC().Format(time.RFC3339))
			if err != nil {
				m.errMsg = err.Error()
				m.starting = false
				return m, nil
			}
			m.activeID = log.ID
			m.activeLbl = log.Label
			m.status = "Timer started"
			m.starting = false
			_ = m.reload()
			return m, tickCmd()
		}
	}

	var cmd tea.Cmd
	m.label, cmd = m.label.Update(msg)
	return m, cmd
}

func (m *Module) View() string {
	if m.starting {
		return fmt.Sprintf("Start Timer (enter to start, esc to cancel):\n\n%s", m.label.View())
	}

	status := "No active timer"
	if m.activeID != 0 {
		if m.activeLbl != "" {
			status = fmt.Sprintf("Active: %s (id %d)", m.activeLbl, m.activeID)
		} else {
			status = fmt.Sprintf("Active timer ID: %d", m.activeID)
		}
	}

	lines := []string{"Tiny Task Timer", status, ""}
	for i, log := range m.logs {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		startText := formatTimestamp(log.Start)
		endText := formatTimestamp(log.End)
		elapsed := formatElapsed(log.Start, log.End)
		line := fmt.Sprintf("%s%s | %s -> %s (%s)", prefix, log.Label, startText, endText, elapsed)
		lines = append(lines, line)
	}
	if len(m.logs) == 0 {
		lines = append(lines, "  (empty) Press s to start a timer.")
	}
	if len(m.logs) > 0 {
		lines = append(lines, fmt.Sprintf("  Selected: %d of %d (use n/p)", m.cursor+1, len(m.logs)))
	}

	if m.errMsg != "" {
		lines = append(lines, "", "Error: "+m.errMsg)
	}
	if m.status != "" {
		lines = append(lines, "", "Status: "+m.status)
	}

	return joinLines(lines)
}

func (m *Module) reload() error {
	logs, err := m.app.Stores.Time.List(context.Background(), 10)
	if err != nil {
		m.errMsg = err.Error()
		return err
	}
	m.logs = logs
	m.activeID = 0
	m.activeLbl = ""
	for _, log := range logs {
		if log.End == "" {
			m.activeID = log.ID
			m.activeLbl = log.Label
			break
		}
	}
	if m.cursor >= len(m.logs) {
		m.cursor = len(m.logs) - 1
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
	return nil
}

func (m *Module) stopTimer() {
	if err := m.app.Stores.Time.Stop(context.Background(), m.activeID, time.Now().UTC().Format(time.RFC3339)); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.activeID = 0
	m.activeLbl = ""
	m.status = "Timer stopped"
	_ = m.reload()
}

func (m *Module) deleteSelected() {
	if len(m.logs) == 0 {
		return
	}
	id := m.logs[m.cursor].ID
	if err := m.app.Stores.Time.Delete(context.Background(), id); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.status = "Deleted timer log"
	_ = m.reload()
}

func (m *Module) clearAll() {
	if err := m.app.Stores.Time.Clear(context.Background()); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.confirm = false
	m.activeID = 0
	m.activeLbl = ""
	m.status = "Cleared all timer logs"
	_ = m.reload()
}

func joinLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	out := lines[0]
	for i := 1; i < len(lines); i++ {
		out += "\n" + lines[i]
	}
	return out
}

func formatElapsed(start, end string) string {
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return "unknown"
	}
	var endTime time.Time
	if end == "" {
		endTime = time.Now().UTC()
	} else {
		endTime, err = time.Parse(time.RFC3339, end)
		if err != nil {
			return "unknown"
		}
	}
	if endTime.Before(startTime) {
		return "0s"
	}
	return formatDuration(endTime.Sub(startTime))
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%dh%dm%ds", h, m, s)
}

func formatTimestamp(value string) string {
	if value == "" {
		return "..."
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return value
	}
	return t.Local().Format("Jan 2 3:04:05 PM MST")
}

func tickCmd() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
