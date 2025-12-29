package footer

import tea "github.com/charmbracelet/bubbletea"

var _ tea.Model = (*Model)(nil)

type Model struct {
}


func (instance *Model) Init() tea.Cmd {
	return nil
}


func (instance *Model) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return instance , nil
}


func (instance *Model) View() string {
	return ""
}