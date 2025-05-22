package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Constants for UI and Mandelbrot parameters
const (
	MoveStep        = 0.1
	WidthAdjustment = 2
	MenuPadding     = 3
)

// UIConfig holds styling and layout configuration
type UIConfig struct {
	MenuWidth     int
	BorderColor   lipgloss.Color
	AccentColor   lipgloss.Color
	TextColor     lipgloss.Color
	HeaderColor   lipgloss.Color
	WhiteColor    lipgloss.Color
	GrayColor     lipgloss.Color
	DisabledColor lipgloss.Color
	ErrorColor    lipgloss.Color
}

var defaultUIConfig = UIConfig{
	MenuWidth:     30,
	BorderColor:   lipgloss.Color("63"),
	AccentColor:   lipgloss.Color("39"),
	TextColor:     lipgloss.Color("250"),
	HeaderColor:   lipgloss.Color("67"),
	WhiteColor:    lipgloss.Color("255"),
	GrayColor:     lipgloss.Color("245"),
	DisabledColor: lipgloss.Color("237"),
	ErrorColor:    lipgloss.Color("196"), // Red
}

// Styling for the UI
var (
	labelStyle = lipgloss.NewStyle().
			Foreground(defaultUIConfig.WhiteColor).
			PaddingLeft(1).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(defaultUIConfig.GrayColor)

	disabledStyle = lipgloss.NewStyle().
			Foreground(defaultUIConfig.DisabledColor).
			PaddingLeft(1).
			Bold(true).
			Strikethrough(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(defaultUIConfig.TextColor)

	headerStyle = lipgloss.NewStyle().
			Foreground(defaultUIConfig.HeaderColor).
			Bold(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(defaultUIConfig.BorderColor)

	panelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(defaultUIConfig.BorderColor).
			Padding(0, 1).
			Width(defaultUIConfig.MenuWidth)

	mandelbrotStyle = lipgloss.NewStyle()

	errorStyle = lipgloss.NewStyle().
			Foreground(defaultUIConfig.ErrorColor)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0"))
)

// styleControlLine styles a single control help line.
func styleControlLine(line string, disabled bool) string {
	parts := strings.SplitN(line, ": ", 2)
	if len(parts) != 2 {
		return valueStyle.Render(line)
	}
	label := labelStyle.Render(parts[0] + ": ")
	value := parts[1]
	if disabled {
		return disabledStyle.Render(parts[0] + ": " + value)
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, label, valueStyle.Render(value))
}
