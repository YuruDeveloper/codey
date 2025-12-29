package styles

import "github.com/charmbracelet/lipgloss"

type Components struct {

	H1 lipgloss.Style
	H2 lipgloss.Style
	Body lipgloss.Style
	BodyDim lipgloss.Style
	Muted lipgloss.Style

	// Box
	Box lipgloss.Style
	BoxFocused lipgloss.Style

	// Input
	Input lipgloss.Style
	InputFocused lipgloss.Style

	// Select
	Select lipgloss.Style
	ListItem lipgloss.Style
	ListItemSelected lipgloss.Style
	ListItemFocused lipgloss.Style

	// Message
	UserMessage lipgloss.Style
	AssistantMessage lipgloss.Style
	SystemMessage lipgloss.Style

	// Status
	StatusBar lipgloss.Style
	ShortcutHint lipgloss.Style

	// Tool
	ToolPending lipgloss.Style
	ToolSuccess lipgloss.Style
	ToolError lipgloss.Style
	ToolDefault lipgloss.Style
}

var DefaultComponents Components = Components{

	H1: lipgloss.NewStyle().Bold(true).Foreground(DefaultColors.Text),
	H2: lipgloss.NewStyle().Bold(true).Foreground(DefaultColors.TextDim),
	Body: lipgloss.NewStyle().Foreground(DefaultColors.Text),
	BodyDim: lipgloss.NewStyle().Foreground(DefaultColors.TextDim),
	Muted: lipgloss.NewStyle().Foreground(DefaultColors.TextMuted),

	Box: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(DefaultColors.Border).PaddingLeft(1).PaddingRight(1),
	BoxFocused: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(DefaultColors.Focus).PaddingLeft(1).PaddingRight(1),

	Input: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(DefaultColors.Border).PaddingLeft(1).PaddingRight(2),
	InputFocused: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(DefaultColors.Focus).PaddingLeft(1).PaddingRight(2),

	Select: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(DefaultColors.Focus),
	ListItem: lipgloss.NewStyle().PaddingLeft(2),
	ListItemSelected: lipgloss.NewStyle().PaddingLeft(2).Foreground(DefaultColors.Focus),
	ListItemFocused: lipgloss.NewStyle().PaddingLeft(0).Foreground(DefaultColors.Focus),

	UserMessage: lipgloss.NewStyle().Foreground(DefaultColors.Text),
	AssistantMessage: lipgloss.NewStyle().Foreground(DefaultColors.Text),
	SystemMessage: lipgloss.NewStyle().Foreground(DefaultColors.TextDim),

	StatusBar: lipgloss.NewStyle().Foreground(DefaultColors.TextDim),
	ShortcutHint: lipgloss.NewStyle().Foreground(DefaultColors.TextMuted),

	ToolPending: lipgloss.NewStyle().Foreground(DefaultColors.Pending),
	ToolSuccess: lipgloss.NewStyle().Foreground(DefaultColors.Success),
	ToolError: lipgloss.NewStyle().Foreground(DefaultColors.Error),
	ToolDefault: lipgloss.NewStyle().Foreground(DefaultColors.Warning),
}