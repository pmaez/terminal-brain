package core

import tea "github.com/charmbracelet/bubbletea"

type Module interface {
	Name() string
	Init(*AppState) error
	Update(tea.Msg) (Module, tea.Cmd)
	View() string
}
