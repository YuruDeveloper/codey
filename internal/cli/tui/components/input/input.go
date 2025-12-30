package input

import (
	"github.com/YuruDeveloper/codey/internal/cli/tui/styles"
	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = (*Model)(nil)

type Model struct {
	textArea textarea.Model
	width int
	maxHeight int
	focused bool
}

var SubmitKey key.Binding = key.NewBinding(key.WithKeys("enter"))

func New() *Model {
	textarea := textarea.New()
	textarea.Placeholder = ""
	textarea.ShowLineNumbers = false
	textarea.SetHeight(1)
	textarea.FocusedStyle.Base = styles.DefaultComponents.Input
	textarea.BlurredStyle.Base = styles.DefaultComponents.Input
	textarea.Focus()
	textarea.SetPromptFunc(2, func (index int) string {
		if index == 0 {
			return styles.DefaultSymbols.Pointer + " " 
		}
		return "  "
	})
	return &Model{
		textArea: textarea,
		maxHeight: 5,
		focused: true,
	}
}

func (instance *Model) Init() tea.Cmd {
	return textarea.Blink
}

func (instance *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch message := message.(type) {
		case tea.WindowSizeMsg:
			instance.width = message.Width
			instance.textArea.SetWidth(message.Width -4)
		case tea.KeyMsg:
			if key.Matches(message,SubmitKey) {
				cmd := func() tea.Msg {
							return types.SubmitInput {
								Text: instance.textArea.Value(),
							}
						}	 
				cmds = append(cmds, cmd)
			} else {
				var cmd tea.Cmd
				instance.textArea , cmd = instance.textArea.Update(message)
				cmds = append(cmds, cmd)
			}
		case types.ResetCommand:
			instance.textArea.Reset()
	}

	contentLen := instance.textArea.Length()
	width := instance.textArea.Width()

	if width > 0 {
		height := (contentLen+1)/width + 1
		if height > instance.maxHeight {
			height = instance.maxHeight
		}
		if height < 1 {
			height = 1
		}
		instance.textArea.SetHeight(height)
	}
	return instance, tea.Batch(cmds...)
}

func (instance *Model) View() string {
	return instance.textArea.View()
}
