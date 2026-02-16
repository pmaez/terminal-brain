package stack

import (
	"context"
	"fmt"

	"terminal-brain/internal/core"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

type Module struct {
	app    *core.AppState
	tasks  []core.StackTask
	cursor int
	input  textinput.Model
	adding bool
	errMsg string
	status string
}

func New() *Module {
	input := textinput.New()
	input.Placeholder = "New task"
	input.Focus()
	return &Module{input: input}
}

func (m *Module) Name() string {
	return "Stack"
}

func (m *Module) Help() string {
	return "Stack: a add | d delete | j/k move | J/K reorder"
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
			m.input.CursorEnd()
			return m, textinput.Blink
		case "d":
			m.deleteSelected()
			return m, nil
		case "j", "down":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
			return m, nil
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "J":
			m.moveSelected(1)
			return m, nil
		case "K":
			m.moveSelected(-1)
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
		case "enter":
			value := m.input.Value()
			if value != "" {
				if _, err := m.app.Stores.Stack.Add(context.Background(), value); err != nil {
					m.errMsg = err.Error()
				} else {
					m.status = "Added task"
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
		return fmt.Sprintf("Add task (enter to save, esc to cancel):\n\n%s", m.input.View())
	}

	lines := []string{fmt.Sprintf("Stack (%d items)", len(m.tasks))}
	for i, task := range m.tasks {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		lines = append(lines, fmt.Sprintf("%s%s", cursor, task.Title))
	}
	if len(m.tasks) == 0 {
		lines = append(lines, "  (empty) Press a to add a task.")
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
	tasks, err := m.app.Stores.Stack.List(context.Background())
	if err != nil {
		m.errMsg = err.Error()
		return err
	}
	m.tasks = tasks
	if m.cursor >= len(m.tasks) {
		m.cursor = len(m.tasks) - 1
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
	return nil
}

func (m *Module) deleteSelected() {
	if len(m.tasks) == 0 {
		return
	}
	id := m.tasks[m.cursor].ID
	if err := m.app.Stores.Stack.Delete(context.Background(), id); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.status = "Deleted task"
	_ = m.reload()
}

func (m *Module) moveSelected(delta int) {
	if len(m.tasks) == 0 {
		return
	}
	target := m.cursor + delta
	if target < 0 || target >= len(m.tasks) {
		return
	}
	m.tasks[m.cursor], m.tasks[target] = m.tasks[target], m.tasks[m.cursor]
	m.cursor = target
	ids := make([]int64, 0, len(m.tasks))
	for _, task := range m.tasks {
		ids = append(ids, task.ID)
	}
	if err := m.app.Stores.Stack.UpdateOrder(context.Background(), ids); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.status = "Reordered task"
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
