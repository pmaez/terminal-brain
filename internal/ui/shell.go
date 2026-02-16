package ui

import (
	"fmt"
	"strings"

	"terminal-brain/internal/core"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ShellModel struct {
	app    *core.AppState
	cursor int
	focusSidebar bool
}

type helpProvider interface {
	Help() string
}

type inputBlocker interface {
	InputActive() bool
}

func NewShellModel(app *core.AppState) *ShellModel {
	return &ShellModel{app: app, focusSidebar: true}
}

func (m *ShellModel) Init() tea.Cmd {
	return nil
}

func (m *ShellModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.focusSidebar = true
			return m, nil
		}
		if m.shouldBypassGlobalKeys() && msg.String() != "esc" {
			break
		}
		if !m.focusSidebar {
			break
		}
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down", "j":
			if m.cursor < len(m.app.Order)-1 {
				m.cursor++
			}
			return m, nil
		case "enter":
			if len(m.app.Order) == 0 {
				return m, nil
			}
			selected := m.app.Order[m.cursor]
			if selected != "" && selected != m.app.Active {
				m.app.Active = selected
				m.focusSidebar = false
				return m, nil
			}
			if selected == m.app.Active {
				m.focusSidebar = false
				return m, nil
			}
		}
	}

	active := m.activeModule()
	if active == nil {
		return m, nil
	}

	updated, cmd := active.Update(msg)
	m.app.Modules[m.app.Active] = updated
	return m, cmd
}

func (m *ShellModel) View() string {
	if len(m.app.Order) == 0 {
		return "No modules registered."
	}

	sidebarWidth := 26
	lines := make([]string, 0, len(m.app.Order))
	activeStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("70"))
	selectedStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("33"))

	for i, key := range m.app.Order {
		name := key
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		if key == m.app.Active {
			name = "* " + name
		} else {
			name = "  " + name
		}
		if key == m.app.Active {
			name = activeStyle.Render(name)
		}
		if i == m.cursor {
			name = selectedStyle.Render(name)
		}
		lines = append(lines, prefix+name)
	}

	sidebarContent := lipgloss.NewStyle().Width(sidebarWidth - 2).Padding(0, 1).Render(strings.Join(lines, "\n"))
	mainContent := lipgloss.NewStyle().Padding(0, 1).Render(m.activeView())

	sidebar := panel("Modules", sidebarContent, sidebarWidth, false)
	main := panel(m.app.Active, mainContent, 0, false)
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)

	footer := panel("Commands", m.footerLine(lipgloss.Width(body)-2), 0, false)
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("70")).
		Width(lipgloss.Width(body)).
		Align(lipgloss.Center).
		Render("Terminal Brain")

	return lipgloss.JoinVertical(lipgloss.Left, title, body, footer)
}

func (m *ShellModel) activeModule() core.Module {
	if m.app.Active == "" {
		return nil
	}
	return m.app.Modules[m.app.Active]
}

func (m *ShellModel) activeView() string {
	active := m.activeModule()
	if active == nil {
		return ""
	}
	return active.View()
}

func (m *ShellModel) shouldBypassGlobalKeys() bool {
	if m.focusSidebar {
		return false
	}
	active := m.activeModule()
	if active == nil {
		return false
	}
	if blocker, ok := active.(inputBlocker); ok {
		return blocker.InputActive()
	}
	return false
}

func (m *ShellModel) footerLine(width int) string {
	focus := "sidebar"
	if !m.focusSidebar {
		focus = "content"
	}
	left := "Global: esc to sidebar (" + focus + ") | q quit | up/down move | enter activate"
	right := ""
	active := m.activeModule()
	if active != nil {
		if helper, ok := active.(helpProvider); ok {
			right = helper.Help()
		}
	}
	if right == "" || width <= 0 {
		return left
	}
	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 2 {
		return left + "\n" + right
	}
	return left + strings.Repeat(" ", gap) + right
}

func panel(title, content string, width int, highlight bool) string {
	borderColor := lipgloss.Color("240")
	titleColor := lipgloss.Color("70")
	if highlight {
		borderColor = lipgloss.Color("118")
		titleColor = lipgloss.Color("118")
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)
	if width > 0 {
		style = style.Width(width)
	}
	header := lipgloss.NewStyle().Bold(true).Foreground(titleColor).Render(fmt.Sprintf("[%s]", title))
	return style.Render(lipgloss.JoinVertical(lipgloss.Left, header, content))
}

func Run(model tea.Model) error {
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
