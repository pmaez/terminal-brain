package cmdmem

import (
	"context"
	"fmt"

	"terminal-brain/internal/core"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

type Module struct {
	app       *core.AppState
	commands  []core.CommandEntry
	cursor    int
	search    textinput.Model
	addCmd    textinput.Model
	addDesc   textinput.Model
	adding    bool
	errMsg    string
	status    string
}

func New() *Module {
	search := textinput.New()
	search.Placeholder = "Search commands"
	search.Prompt = "> "
	search.Focus()

	addCmd := textinput.New()
	addCmd.Placeholder = "Command"
	addDesc := textinput.New()
	addDesc.Placeholder = "Description"

	return &Module{search: search, addCmd: addCmd, addDesc: addDesc}
}

func (m *Module) Name() string {
	return "CmdMem"
}

func (m *Module) Help() string {
	return "CmdMem: type to search | a add | enter mark used"
}

func (m *Module) InputActive() bool {
	return true
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
			m.addCmd.SetValue("")
			m.addDesc.SetValue("")
			m.addCmd.Focus()
			return m, nil
		case "enter":
			m.touchSelected()
			return m, nil
		case "down", "j":
			if m.cursor < len(m.commands)-1 {
				m.cursor++
			}
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.search, cmd = m.search.Update(msg)
	_ = m.reload()
	return m, cmd
}

func (m *Module) updateAdd(msg tea.Msg) (core.Module, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.adding = false
			return m, nil
		case "enter":
			if m.addCmd.Focused() {
				m.addCmd.Blur()
				m.addDesc.Focus()
				return m, nil
			}
			cmdValue := m.addCmd.Value()
			descValue := m.addDesc.Value()
			if cmdValue != "" {
				if _, err := m.app.Stores.Command.Add(context.Background(), cmdValue, descValue); err != nil {
					m.errMsg = err.Error()
				} else {
					m.status = "Saved command"
				}
				_ = m.reload()
			}
			m.adding = false
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.addCmd.Focused() {
		m.addCmd, cmd = m.addCmd.Update(msg)
	} else {
		m.addDesc, cmd = m.addDesc.Update(msg)
	}
	return m, cmd
}

func (m *Module) View() string {
	if m.adding {
		return fmt.Sprintf("Add Command (enter to advance/save, esc to cancel):\n\nCommand: %s\nDescription: %s", m.addCmd.View(), m.addDesc.View())
	}

	lines := []string{"Command Memory", "", m.search.View()}
	for i, cmd := range m.commands {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		lines = append(lines, fmt.Sprintf("%s%s — %s", cursor, cmd.Command, cmd.Description))
	}
	if len(m.commands) == 0 {
		lines = append(lines, "  (no matches) Try a different search or press a to add.")
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
	commands, err := m.app.Stores.Command.Search(context.Background(), m.search.Value(), 20)
	if err != nil {
		m.errMsg = err.Error()
		return err
	}
	m.commands = commands
	if m.cursor >= len(m.commands) {
		m.cursor = len(m.commands) - 1
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
	return nil
}

func (m *Module) touchSelected() {
	if len(m.commands) == 0 {
		return
	}
	id := m.commands[m.cursor].ID
	if err := m.app.Stores.Command.Touch(context.Background(), id); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.status = "Marked command used"
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
