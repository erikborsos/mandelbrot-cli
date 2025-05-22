package tui

import (
	"fmt"
	"mandel-cli/kitty"
	"mandel-cli/mandelbrot"
	"mandel-cli/utils"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyAction defines possible user actions
type KeyAction string

const (
	MoveLeft     KeyAction = "move_left"
	MoveRight    KeyAction = "move_right"
	MoveUp       KeyAction = "move_up"
	MoveDown     KeyAction = "move_down"
	ZoomIn       KeyAction = "zoom_in"
	ZoomOut      KeyAction = "zoom_out"
	CycleColor   KeyAction = "cycle_color"
	ToggleSmooth KeyAction = "toggle_smooth"
	IncreaseIter KeyAction = "increase_iter"
	DecreaseIter KeyAction = "decrease_iter"
	Reset        KeyAction = "reset"
	ToggleImg    KeyAction = "toggle_img"
	Quit         KeyAction = "quit"
	Hide         KeyAction = "hide"
)

// KeyHandler defines a function that modifies the model
type KeyHandler func(*Model)

// Key bindings and handlers
var keyBindings = map[KeyAction][]string{
	MoveLeft:     {"h", "left"},
	MoveRight:    {"l", "right"},
	MoveUp:       {"k", "up"},
	MoveDown:     {"j", "down"},
	ZoomIn:       {"+"},
	ZoomOut:      {"-"},
	CycleColor:   {"c"},
	ToggleSmooth: {"s"},
	IncreaseIter: {"i"},
	DecreaseIter: {"d"},
	Reset:        {"r"},
	ToggleImg:    {"t"},
	Quit:         {"q", "ctrl+c"},
	Hide:         {"m"},
}

var keyHandlers = map[KeyAction]KeyHandler{
	MoveLeft:     func(m *Model) { m.params.Move(-MoveStep, 0); m.paramsChanged = true },
	MoveRight:    func(m *Model) { m.params.Move(MoveStep, 0); m.paramsChanged = true },
	MoveUp:       func(m *Model) { m.params.Move(0, -MoveStep); m.paramsChanged = true },
	MoveDown:     func(m *Model) { m.params.Move(0, MoveStep); m.paramsChanged = true },
	ZoomIn:       func(m *Model) { m.params.ZoomIn(); m.paramsChanged = true },
	ZoomOut:      func(m *Model) { m.params.ZoomOut(); m.paramsChanged = true },
	CycleColor:   func(m *Model) { m.params.CycleColor(); m.paramsChanged = true },
	ToggleSmooth: func(m *Model) { m.params.ToggleSmooth(); m.paramsChanged = true },
	IncreaseIter: func(m *Model) { m.params.IncreaseIterations(); m.paramsChanged = true },
	DecreaseIter: func(m *Model) { m.params.DecreaseIterations(); m.paramsChanged = true },
	Reset:        func(m *Model) { m.Reset(); m.paramsChanged = true },
	ToggleImg:    func(m *Model) { m.toggleDisplayImg() },
	Hide:         func(m *Model) { m.toggleHideMenu() },
}

// Model represents the state of the Mandelbrot TUI application.
type Model struct {
	params        mandelbrot.MandelbrotParams // Parameters for Mandelbrot rendering
	text          string                      // Text representation of the Mandelbrot set
	image         string                      // Kitty terminal image representation
	width         int                         // Terminal width
	height        int                         // Terminal height
	displayImg    bool                        // Whether to display image or text
	paramsChanged bool                        // Whether parameters have changed
	errorMsg      string                      // Error message for UI display
	hideMenu      bool
}

// InitModel initializes a new Model with default Mandelbrot parameters.
func InitModel() Model {
	return Model{
		params:   mandelbrot.InitialMandelbrotParams(),
		hideMenu: false,
	}
}

// setupUIConfig initializes UI-related configuration.
func setupUIConfig() {
	infoReplacer = *utils.NewReplacer(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Center Re: "), valueStyle.Render(":CENTER_RE:")),
			lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Center Im: "), valueStyle.Render(":CENTER_IM:")),
			lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Zoom: "), valueStyle.Render(":ZOOM:")),
			lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Iterations: "), valueStyle.Render(":ITER:")),
			lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Color: "), valueStyle.Render(":COLOR:")),
			lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Smooth: "), valueStyle.Render(":SMOOTH:")),
		))

	var helpText = []string{
		"h/j/k/l or arrows: Move",
		"+/-: Zoom in/out",
		"c: Cycle color scheme",
		"s: Toggle smooth coloring",
		"i/d: +/- max iterations",
		"r: Reset to default",
		"t: Toggle image/text",
		"h: Hide menu",
		"q: Quit",
	}

	generateControls := func(disableNonToggle bool) string {
		controlsArr := make([]string, len(helpText))
		for i, line := range helpText {
			controlsArr[i] = styleControlLine(line, disableNonToggle && !strings.HasPrefix(line, "t:") && !strings.HasPrefix(line, "q:"))
		}
		return lipgloss.JoinVertical(lipgloss.Left, controlsArr...)
	}

	controls = generateControls(false)
	controlsDisabled = generateControls(true)
}

// Init initializes the model and sets up UI configuration.
func (m Model) Init() tea.Cmd {
	setupUIConfig()
	return nil
}

// Reset restores the model to default parameters while preserving window size.
func (m *Model) Reset() {
	currentWidth := m.params.Width
	currentHeight := m.params.Height
	m.params.Reset()
	m.params.Width = currentWidth
	m.params.Height = currentHeight
	m.errorMsg = ""
}

// toggleDisplayImg switches between text and image display modes.
func (m *Model) toggleDisplayImg() {
	m.displayImg = !m.displayImg
	m.errorMsg = ""
	if m.displayImg {
		image, err := mandelbrot.GenerateMandelbrotImage(m.params)
		if err != nil {
			m.errorMsg = fmt.Sprintf("Error generating image: %v", err)
			m.displayImg = false
			return
		}
		m.image, err = kitty.Kitty(image, m.params.Width*2, m.params.Height)
		if err != nil {
			m.errorMsg = fmt.Sprintf("Error rendering kitty image: %v", err)
			m.displayImg = false
			return
		}
		m.text = ""
	} else {
		kitty.KittyClearImages()
		m.image = ""
		m.paramsChanged = true
	}
}

func (m *Model) toggleHideMenu() {
	m.hideMenu = !m.hideMenu
	if m.hideMenu {
		m.params.Width = m.width / 2
	} else {
		m.params.Width = (m.width-defaultUIConfig.MenuWidth)/2 - WidthAdjustment
	}

	m.paramsChanged = true
}

// Update handles messages and updates the model state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		for action, keys := range keyBindings {
			for _, k := range keys {
				if k == key {
					if action == Quit {
						return m, tea.Quit
					}
					if !m.displayImg || action == ToggleImg {
						keyHandlers[action](&m)
					}
					break
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.hideMenu {
			m.params.Width = m.width / 2
		} else {
			m.params.Width = (m.width-defaultUIConfig.MenuWidth)/2 - WidthAdjustment
		}
		m.params.Height = m.height
		m.paramsChanged = true
	}
	if !m.displayImg && m.paramsChanged {
		m.text = mandelbrot.BufferToString(mandelbrot.GenerateMandelbrotText(m.params))
		m.paramsChanged = false
	}
	return m, nil
}

// View renders the UI with a vertical split between Mandelbrot and menu.
func (m Model) View() string {
	if !m.hideMenu {
		info := infoReplacer
		infoStr := info.
			Replace(":CENTER_RE:", fmt.Sprintf("%.9f", m.params.CenterRe)).
			Replace(":CENTER_IM:", fmt.Sprintf("%.9f", m.params.CenterIm)).
			Replace(":ZOOM:", fmt.Sprintf("%.9f", m.params.ZoomFactor)).
			Replace(":ITER:", fmt.Sprintf("%d", m.params.MaxIter)).
			Replace(":COLOR:", mandelbrot.ColorNames[m.params.ColorMode]).
			Replace(":SMOOTH:", fmt.Sprintf("%v", m.params.Smooth)).
			String()

		errorStr := ""
		if m.errorMsg != "" {
			errorStr = errorStyle.Render("Error: " + m.errorMsg)
		}

		menuContent := lipgloss.JoinVertical(
			lipgloss.Left,
			headerStyle.Render("Parameters"),
			infoStr,
			headerStyle.Render("Controls"),
			helpStyle.Render(utils.Ternary(m.displayImg, controlsDisabled, controls)),
			errorStr,
		)

		mandelbrotPanel := mandelbrotStyle.
			Width(m.width - defaultUIConfig.MenuWidth - MenuPadding).
			Height(m.height - 2).
			Render(utils.Ternary(m.displayImg, utils.PadEmptyLines(m.image, m.params.Height), m.text))

		menuPanel := panelStyle.Render(menuContent)

		return lipgloss.JoinHorizontal(
			lipgloss.Center,
			mandelbrotPanel,
			lipgloss.NewStyle().Width(1).Render(""),
			menuPanel,
		)
	} else {
		return utils.Ternary(m.displayImg, utils.PadEmptyLines(m.image, m.params.Height), m.text)
	}
}
