package dump

import (
	"context"
	"fmt"

	"terminal-brain/internal/core"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textarea"
)

type Module struct {
	app     *core.AppState
	entries []core.DumpEntry
	input   textarea.Model
	adding  bool
	errMsg  string
	status  string
}

func New() *Module {
	input := textarea.New()
	input.Placeholder = "Dump your thoughts..."
	input.SetHeight(5)
	return &Module{input: input}
}

func (m *Module) Name() string {
	return "Dump"
}

func (m *Module) Help() string {
	return "Dump: a add | ctrl+s save | esc cancel"
}

func (m *Module) InputActive() bool {
	return m.adding
}

func (m *Module) Init(app *core.AppState) error {
	m.app = app
	return m.reload()
}

func (m *Module) Update(msg tea.Msg) (core.Module, tea.Cmd) {
	if m.adding {
		return m.updateAdd(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.adding = true
			m.input.SetValue("")
			m.input.Focus()
			return m, nil
		case "r":
			_ = m.reload()
			return m, nil
		}
	}

	return m, nil
}

func (m *Module) updateAdd(msg tea.Msg) (core.Module, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.adding = false
			return m, nil
		case "ctrl+s":
			value := m.input.Value()
			if value != "" {
				if _, err := m.app.Stores.Dump.Add(context.Background(), value); err != nil {
					m.errMsg = err.Error()
				} else {
					m.status = "Saved dump entry"
				}
				_ = m.reload()
			}
			m.adding = false
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Module) View() string {
	if m.adding {
		return fmt.Sprintf("Brain Dump (ctrl+s to save, esc to cancel):\n\n%s", m.input.View())
	}

	lines := []string{"Brain Dump (most recent first)"}
	for _, entry := range m.entries {
		snippet := entry.Body
		if len(snippet) > 60 {
			snippet = snippet[:60] + "..."
		}
		lines = append(lines, fmt.Sprintf("- %s", snippet))
	}
	if len(m.entries) == 0 {
		lines = append(lines, "  (empty) Press a to add a dump entry.")
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
	entries, err := m.app.Stores.Dump.List(context.Background(), 10)
	if err != nil {
		m.errMsg = err.Error()
		return err
	}
	m.entries = entries
	return nil
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
