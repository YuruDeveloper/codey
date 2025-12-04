package tui

import (
	"github.com/rivo/tview"
)

type tui struct {
	app *tview.Application
}

func Init() *tui{
	
	app := tview.NewApplication()
	return &tui {
		app: app,
	}
}