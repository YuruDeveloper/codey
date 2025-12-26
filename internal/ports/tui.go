package ports

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

type AppTui interface {
	Run(ctx context.Context) error
	SetRouter(router AppRouter)
	Program() *tea.Program
	Send(msg tea.Msg)
	Quit()
}

