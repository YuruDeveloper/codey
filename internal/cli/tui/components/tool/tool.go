package tool

import (
	"fmt"

	"github.com/YuruDeveloper/codey/internal/cli/tui/styles"
	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = (*Model)(nil)

type Model struct {
	cursor   cursor.Model
	toolName string
	path     string
	info     string
	status types.ToolStatus
}

func New(toolName, path string) Model{
	cursor := cursor.New()
	cursor.SetChar(styles.DefaultSymbols.Bullet)
	cursor.Style = styles.DefaultComponents.ToolPending
	return Model{
		toolName: toolName,
		path: path,
		status: types.ToolPending,
	}
}

func (instance *Model) Init() tea.Cmd {
	return instance.cursor.BlinkCmd()
}


func (instance *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if message , ok := message.(types.UpdateToolStatus) ; ok {
		instance.status = message.Status
		instance.info = message.Info
		instance.cursor.SetMode(cursor.CursorStatic)
		switch message.Status {
			case types.ToolSuccess:
				instance.cursor.Style = styles.DefaultComponents.ToolSuccess
			case types.ToolError:
				instance.cursor.Style = styles.DefaultComponents.ToolError
			default:
				instance.cursor.Style = styles.DefaultComponents.ToolDefault
		}
		return instance , nil
	}
	instance.cursor, cmd = instance.cursor.Update(message)
	return instance, cmd 
}


func (instance *Model) View() string {
	view := fmt.Sprintf("%s %s(%s)",instance.cursor.View(),instance.toolName,instance.path)
	if instance.status != types.ToolPending {
		view = lipgloss.JoinVertical(lipgloss.Left,view,fmt.Sprintf("⎿  %s",instance.info))
	}
	return view
}
