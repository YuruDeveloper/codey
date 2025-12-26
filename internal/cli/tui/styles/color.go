package styles

import (
	"github.com/charmbracelet/lipgloss"
)

type Colors struct {
	Text lipgloss.ANSIColor
	TextDimmation lipgloss.ANSIColor
	TextMuted lipgloss.ANSIColor
	
	Border lipgloss.ANSIColor

	Success lipgloss.ANSIColor
	Error lipgloss.ANSIColor
	Warning lipgloss.ANSIColor
	Info lipgloss.ANSIColor

	Focus lipgloss.ANSIColor
	Highlight lipgloss.ANSIColor
	Pending lipgloss.ANSIColor
}

// Default  Colors Palette
var DefaultColors Colors = Colors{
	Text: lipgloss.ANSIColor(255),
	TextDimmation: lipgloss.ANSIColor(245),
	TextMuted: lipgloss.ANSIColor(240),
	Border: lipgloss.ANSIColor(8),

	Success: lipgloss.ANSIColor(10),
	Error: lipgloss.ANSIColor(9),
	Warning: lipgloss.ANSIColor(11),
	Info: lipgloss.ANSIColor(12),

	Focus: lipgloss.ANSIColor(32),
	Highlight: lipgloss.ANSIColor(14),
	Pending: lipgloss.ANSIColor(8),
}