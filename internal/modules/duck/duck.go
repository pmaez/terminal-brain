package duck

import (
	"context"
	"fmt"

	"terminal-brain/internal/core"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

type Module struct {
	app      *core.AppState
	inputs   []textinput.Model
	step     int
	errMsg   string
	status  string
	sessions []core.DuckSession
	mode     duckMode
	selected int
	editID   int64
}

var questions = []string{
	"What did you expect?",
	"What happened?",
	"What did you try?",
	"What changed recently?",
}

func New() *Module {
	inputs := make([]textinput.Model, len(questions))
	for i := range inputs {
		input := textinput.New()
		input.Placeholder = questions[i]
		inputs[i] = input
	}
	inputs[0].Focus()
	return &Module{inputs: inputs}
}

func (m *Module) Name() string {
	return "Duck"
}

func (m *Module) Help() string {
	switch m.mode {
	case modeList:
		return "Duck: up/down move | enter view | esc back"
	case modeView:
		return "Duck: esc back | e edit | d delete"
	case modeEdit:
		return "Duck: enter next/save | ctrl+s save | esc cancel edit"
	case modeSaved:
		return "Duck: ctrl+n new session | v view"
	default:
		return "Duck: enter next/save | up/down move | esc cancel input | ctrl+r reset | ctrl+g view | d delete"
	}
}

func (m *Module) InputActive() bool {
	return m.mode == modeForm || m.mode == modeEdit
}

func (m *Module) Init(app *core.AppState) error {
	m.app = app
	m.mode = modeForm
	return m.reload()
}

func (m *Module) Update(msg tea.Msg) (core.Module, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode == modeList {
			switch msg.String() {
			case "esc":
				m.reset()
				m.mode = modeForm
				return m, nil
			case "up":
				if m.selected > 0 {
					m.selected--
				}
				return m, nil
			case "down":
				if m.selected < len(m.sessions)-1 {
					m.selected++
				}
				return m, nil
			case "enter":
				if len(m.sessions) > 0 {
					m.mode = modeView
				}
				return m, nil
			}
		}
		if m.mode == modeView {
			switch msg.String() {
			case "esc":
				m.reset()
				m.mode = modeForm
				return m, nil
			case "e":
				if m.selected < len(m.sessions) {
					m.startEdit(m.sessions[m.selected])
				}
				return m, nil
			case "d":
				m.deleteSelected()
				return m, nil
			}
		}
		if m.mode == modeSaved {
			switch msg.String() {
			case "ctrl+n":
				m.reset()
				m.editID = 0
				m.mode = modeForm
				m.status = "New session"
				return m, nil
			case "ctrl+g":
				if len(m.sessions) > 0 {
					m.mode = modeList
				}
				return m, nil
			}
		}
		if m.mode == modeEdit {
			switch msg.String() {
			case "esc":
				m.mode = modeForm
				m.reset()
				return m, nil
			case "ctrl+s":
				m.saveEdit()
				return m, nil
			}
		}
		switch msg.String() {
		case "up":
			if m.step > 0 {
				m.step--
				m.inputs[m.step].Focus()
			}
			return m, nil
		case "down":
			if m.step < len(m.inputs)-1 {
				m.step++
				m.inputs[m.step].Focus()
			}
			return m, nil
		case "esc":
			m.inputs[m.step].SetValue("")
			m.status = "Canceled current input"
			return m, nil
		case "ctrl+r":
			if m.mode == modeForm {
				m.reset()
				m.status = "Reset prompts"
			}
			return m, nil
		case "enter":
			if m.step < len(m.inputs)-1 {
				m.step++
				m.inputs[m.step].Focus()
				return m, nil
			}
			if m.mode == modeEdit {
				m.saveEdit()
				return m, nil
			}
			m.saveNew()
			return m, nil
		case "ctrl+g":
			if (m.mode == modeForm || m.mode == modeSaved) && len(m.sessions) > 0 {
				m.mode = modeList
				return m, nil
			}
		case "d":
			if (m.mode == modeForm || m.mode == modeSaved) && len(m.sessions) > 0 {
				m.deleteSelected()
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	if m.mode == modeForm || m.mode == modeEdit {
		m.inputs[m.step], cmd = m.inputs[m.step].Update(msg)
	}
	return m, cmd
}

func (m *Module) View() string {
	if m.mode == modeSaved {
		lines := []string{
			"Saved duck session.",
			"",
			"Press ctrl+n to start a new session.",
			"Press ctrl+g to view recent sessions.",
		}
		return joinLines(lines)
	}
	if m.mode == modeView && len(m.sessions) > 0 {
		session := m.sessions[m.selected]
		lines := []string{
			"Saved Session",
			"",
			"1. " + session.Q1,
			"2. " + session.Q2,
			"3. " + session.Q3,
			"4. " + session.Q4,
			"",
			"Press e to edit or esc to return.",
		}
		return joinLines(lines)
	}
	if m.mode == modeList && len(m.sessions) > 0 {
		lines := []string{"Recent Sessions", ""}
		for i, session := range m.sessions {
			prefix := "  "
			if i == m.selected {
				prefix = "> "
			}
			snippet := session.Q1
			if len(snippet) > 40 {
				snippet = snippet[:40] + "..."
			}
			lines = append(lines, fmt.Sprintf("%s%s", prefix, snippet))
		}
		lines = append(lines, "", "Press enter to view, esc to return.")
		return joinLines(lines)
	}

	lines := []string{"Rubber Duck Mode (enter to continue, esc cancel input):", ""}
	for i, input := range m.inputs {
		prefix := "  "
		if i == m.step {
			prefix = "> "
		}
		lines = append(lines, fmt.Sprintf("%s%s", prefix, input.View()))
	}
	lines = append(lines, "", "Tip: short answers are fine. You can refine later.")
	if len(m.sessions) > 0 {
		lines = append(lines, "", "Recent sessions:")
		for i, session := range m.sessions {
			snippet := session.Q1
			if len(snippet) > 40 {
				snippet = snippet[:40] + "..."
			}
			prefix := fmt.Sprintf("%d.", i+1)
			lines = append(lines, fmt.Sprintf("%s %s", prefix, snippet))
		}
		lines = append(lines, "Press ctrl+g to view recent sessions.")
	}

	if m.errMsg != "" {
		lines = append(lines, "", "Error: "+m.errMsg)
	}
	if m.status != "" {
		lines = append(lines, "", "Status: "+m.status)
	}

	return joinLines(lines)
}

func (m *Module) saveNew() {
	_, err := m.app.Stores.Duck.Add(
		context.Background(),
		m.inputs[0].Value(),
		m.inputs[1].Value(),
		m.inputs[2].Value(),
		m.inputs[3].Value(),
	)
	if err != nil {
		m.errMsg = err.Error()
		return
	}
	m.status = "Saved duck session"
	_ = m.reload()
	m.mode = modeSaved
}

func (m *Module) saveEdit() {
	if m.editID == 0 {
		return
	}
	if err := m.app.Stores.Duck.Update(
		context.Background(),
		m.editID,
		m.inputs[0].Value(),
		m.inputs[1].Value(),
		m.inputs[2].Value(),
		m.inputs[3].Value(),
	); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.status = "Updated duck session"
	m.mode = modeForm
	m.editID = 0
	_ = m.reload()
	m.reset()
	m.status = "Updated duck session"
}

func (m *Module) startEdit(session core.DuckSession) {
	m.inputs[0].SetValue(session.Q1)
	m.inputs[1].SetValue(session.Q2)
	m.inputs[2].SetValue(session.Q3)
	m.inputs[3].SetValue(session.Q4)
	m.step = 0
	m.inputs[0].Focus()
	m.editID = session.ID
	m.mode = modeEdit
}

func (m *Module) deleteSelected() {
	if len(m.sessions) == 0 {
		return
	}
	id := m.sessions[m.selected].ID
	if err := m.app.Stores.Duck.Delete(context.Background(), id); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.status = "Deleted duck session"
	m.mode = modeForm
	_ = m.reload()
}

func (m *Module) reset() {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}
	m.step = 0
	m.inputs[0].Focus()
	if m.mode == modeEdit {
		m.mode = modeForm
		m.editID = 0
	}
}

func (m *Module) reload() error {
	sessions, err := m.app.Stores.Duck.List(context.Background(), 5)
	if err != nil {
		m.errMsg = err.Error()
		return err
	}
	m.sessions = sessions
	if m.selected >= len(m.sessions) {
		m.selected = len(m.sessions) - 1
		if m.selected < 0 {
			m.selected = 0
		}
	}
	return nil
}

type duckMode int

const (
	modeForm duckMode = iota
	modeList
	modeView
	modeEdit
	modeSaved
)

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
