package tui

import (
	"context"
	"github.com/YuruDeveloper/codey/internal/ports"
	tea "github.com/charmbracelet/bubbletea"
)

var _ ports.AppTui = (*Tui)(nil)

type Tui struct {
	program *tea.Program
	model *Model
	router ports.AppRouter
}

func (instance *Tui) Program() *tea.Program {
	return instance.program
}


func (instance *Tui) Quit() {
	instance.Quit()
}

func (instance *Tui) Run(ctx context.Context) error {
	instance.model = new()
	instance.program = tea.NewProgram(
		instance.model,
	)

	go func() {
		<-ctx.Done()
		instance.Quit()
	}()
	
	_ , err := instance.program.Run()
	return err
}
func (instance *Tui) Send(msg tea.Msg) {
	instance.program.Send(msg)
}

func (instance *Tui) SetRouter(router ports.AppRouter) {
	instance.router = router
}
