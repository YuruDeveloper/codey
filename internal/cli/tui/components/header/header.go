package header

import (
	"fmt"

	"github.com/YuruDeveloper/codey/internal/cli/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = (*Model)(nil)

const character = "█████\n█▛█▜█\n█████\n▘ ▘ ▘\n▘ ▘ ▘\n"

func New(appName, version, model, workingDir string) Model {
	return Model{
		appName: appName,
		version: version,
		model: model,
		workingDir: workingDir,
	}
}

type Model struct {
	appName   string
	version   string
	model string
	workingDir string
	width     int
}

func (instance *Model) Init() tea.Cmd {
	return nil
}

func (instance *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		instance.width = message.Width
	}
	return instance , nil
}

func (instance *Model) View() string {
	nameLine := fmt.Sprintf("%s %s",styles.DefaultComponents.H1.Render(instance.appName),styles.DefaultComponents.BodyDim.Render(instance.version))
	modelLine := styles.DefaultComponents.BodyDim.Render(fmt.Sprintf("Model: %s",instance.model))
	dirLine := styles.DefaultComponents.Muted.Render(instance.workingDir)
	line := lipgloss.JoinVertical(lipgloss.Left,nameLine,modelLine,dirLine)
	return lipgloss.JoinHorizontal(lipgloss.Bottom,character,line)
}
