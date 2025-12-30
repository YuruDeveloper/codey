package tui

import (
	"context"
	"fmt"

	"github.com/YuruDeveloper/codey/internal/cli/tui/components/input"
	"github.com/YuruDeveloper/codey/internal/cli/tui/components/tool"
	"github.com/YuruDeveloper/codey/internal/cli/tui/styles"
	"github.com/YuruDeveloper/codey/internal/ports"
	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

var _ tea.Model = (*Model)(nil)

var ExitKey key.Binding = key.NewBinding(key.WithKeys("ctrl+c"))

type ToolItem struct {
	Model tea.Model
	UUID uuid.UUID
}

func new() *Model {
	return &Model{
		textInput: input.New(),
		selectInput: nil,
		status: types.TuiStateTalk,
		toolsSlice : make([]ToolItem,0,3),
	}
}

type Model struct {
	textInput tea.Model
	selectInput tea.Model
	status types.TuiState
	toolsSlice []ToolItem
	router ports.AppRouter
}

func (instance *Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmd := instance.textInput.Init()
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (instance *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch message := message.(type) {
		case types.StartTool:
			newTool := tool.New(message.Name,message.Path,message.UUID,instance.DeleteToolFromMap)
			instance.toolsSlice = append(instance.toolsSlice, ToolItem{ Model: newTool, UUID: message.UUID})
			cmd := newTool.Init()
			cmds = append(cmds, cmd)
		case types.SubmitInput:
			if len(message.Text) == 0 {
				break
			}
			tea.Println(fmt.Sprintf("%s %s",styles.DefaultSymbols.Cursor,message.Text))
			cmd := func () tea.Msg {
				return types.ResetCommand{}
			}
			cmds = append(cmds, cmd)
			instance.router.SendMessage(context.Background(),[]types.Part{ types.TextPart{ Text: message.Text } })
		case tea.KeyMsg:
			if key.Matches(message,ExitKey) {
				return  instance, tea.Quit
			}
	}
	var cmd tea.Cmd
	for index := range instance.toolsSlice {
		instance.toolsSlice[index].Model , cmd = instance.toolsSlice[index].Model.Update(message)
		cmds = append(cmds, cmd)
	}
	switch instance.status {
		case types.TuiStateTalk:
			var cmd tea.Cmd
			instance.textInput , cmd = instance.textInput.Update(message)
			cmds = append(cmds, cmd)
	}
	return instance, tea.Batch(cmds...)
}


func (instance *Model) View() string {
	var station []string
	for id := range instance.toolsSlice {
		station = append(station, instance.toolsSlice[id].Model.View())
	}
	station = append(station, instance.textInput.View())
	return lipgloss.JoinVertical(lipgloss.Left,station...)
}

func  (instance *Model) DeleteToolFromMap(uuid uuid.UUID) {
	for index := range instance.toolsSlice {
	 	if instance.toolsSlice[index].UUID == uuid { 
			instance.toolsSlice = append(instance.toolsSlice[:index],instance.toolsSlice[index+1:]...)
		}
	}
}

func (instance *Model) Quit() {
	instance.Update(tea.Quit())
}

func (instance *Model) SetRouter(router ports.AppRouter) {
	instance.router = router
}
