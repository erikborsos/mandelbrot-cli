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
	ForceQuit    KeyAction = "force_quit"
	Hide         KeyAction = "hide"
	SelectPreset KeyAction = "select_preset"
	Save         KeyAction = "save"
)

type KeyHandler func(*Model)

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
	Quit:         {"q"},
	ForceQuit:    {"ctrl+c"},
	Hide:         {"m"},
	SelectPreset: {"p"},
	Save:         {"ctrl+s"},
}

var mandelbrotKeyHandlers = map[KeyAction]KeyHandler{
	MoveLeft:     func(m *Model) { m.params.Move(-MoveStep, 0); m.mandelbortModel.paramsChanged = true },
	MoveRight:    func(m *Model) { m.params.Move(MoveStep, 0); m.mandelbortModel.paramsChanged = true },
	MoveUp:       func(m *Model) { m.params.Move(0, -MoveStep); m.mandelbortModel.paramsChanged = true },
	MoveDown:     func(m *Model) { m.params.Move(0, MoveStep); m.mandelbortModel.paramsChanged = true },
	ZoomIn:       func(m *Model) { m.params.ZoomIn(); m.mandelbortModel.paramsChanged = true },
	ZoomOut:      func(m *Model) { m.params.ZoomOut(); m.mandelbortModel.paramsChanged = true },
	CycleColor:   func(m *Model) { m.params.CycleColor(); m.mandelbortModel.paramsChanged = true },
	ToggleSmooth: func(m *Model) { m.params.ToggleSmooth(); m.mandelbortModel.paramsChanged = true },
	IncreaseIter: func(m *Model) { m.params.IncreaseIterations(); m.mandelbortModel.paramsChanged = true },
	DecreaseIter: func(m *Model) { m.params.DecreaseIterations(); m.mandelbortModel.paramsChanged = true },
	Reset:        func(m *Model) { m.Reset(); m.mandelbortModel.paramsChanged = true },
	ToggleImg:    func(m *Model) { m.toggleDisplayImg() },
	Hide:         func(m *Model) { m.toggleHideMenu() },
	SelectPreset: func(m *Model) {
		m.view = PresetsView
		h, v := docStyle.GetFrameSize()
		m.presetsModel.list.SetSize(m.width-h, m.height-v)
	},
	Save: func(m *Model) {
		m.view = PresetsView
		m.saveModel = initSaveModel(m.params)
		m.view = SaveView
	},
}

var infoReplacer utils.ChainReplacer
var controls string
var controlsDisabled string

type MandelbrotModel struct {
	text          string // Text representation of the Mandelbrot set
	image         string // Kitty terminal image representation
	displayImg    bool   // Whether to display image or text
	paramsChanged bool   // Whether parameters have changed
	errorMsg      string // Error message for UI display
	hideMenu      bool   // Wheter menu should be hidden
}

func initMandelbrotModel() MandelbrotModel {
	return MandelbrotModel{
		hideMenu:      false,
		paramsChanged: true,
	}
}

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
		"p: Select preset",
		"m: Hide menu",
		"t: Toggle image/text",
		"q: Quit",
	}

	generateControls := func(disableNonToggle bool) string {
		controlsArr := make([]string, len(helpText))
		for i, line := range helpText {
			controlsArr[i] = styleControlLine(line, disableNonToggle && !strings.HasPrefix(line, "t:") && !strings.HasPrefix(line, "q:") && !strings.HasPrefix(line, "m:"))
		}
		return lipgloss.JoinVertical(lipgloss.Left, controlsArr...)
	}

	controls = generateControls(false)
	controlsDisabled = generateControls(true)
}

func (m *Model) toggleDisplayImg() {
	m.mandelbortModel.displayImg = !m.mandelbortModel.displayImg
	m.mandelbortModel.errorMsg = ""
	if m.mandelbortModel.displayImg {
		image, err := mandelbrot.GenerateMandelbrotImage(m.params)
		if err != nil {
			m.mandelbortModel.errorMsg = fmt.Sprintf("Error generating image: %v", err)
			m.mandelbortModel.displayImg = false
			return
		}
		m.mandelbortModel.image, err = kitty.Kitty(image, m.params.Width*2, m.params.Height)
		if err != nil {
			m.mandelbortModel.errorMsg = fmt.Sprintf("Error rendering kitty image: %v", err)
			m.mandelbortModel.displayImg = false
			return
		}
		m.mandelbortModel.text = ""
	} else {
		kitty.KittyClearImages()
		m.mandelbortModel.image = ""
		m.mandelbortModel.paramsChanged = true
	}
}

func (m *Model) toggleHideMenu() {
	m.mandelbortModel.hideMenu = !m.mandelbortModel.hideMenu
	if m.mandelbortModel.hideMenu {
		m.params.Width = m.width / 2
	} else {
		m.params.Width = (m.width-defaultUIConfig.MenuWidth)/2 - WidthAdjustment
	}

	m.mandelbortModel.paramsChanged = true
}

func (m Model) ViewMandelbrot() string {
	if !m.mandelbortModel.hideMenu {
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
		if m.mandelbortModel.errorMsg != "" {
			errorStr = errorStyle.Render("Error: " + m.mandelbortModel.errorMsg)
		}

		menuContent := lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			headerStyle.Render("Parameters:"),
			lipgloss.NewStyle().Padding(0, 0, 1, 0).Render(infoStr),
			headerStyle.Render("Controls:"),
			helpStyle.Render(utils.Ternary(m.mandelbortModel.displayImg, controlsDisabled, controls)),
			errorStr,
		)

		mandelbrotPanel := mandelbrotStyle.
			Width(m.width - defaultUIConfig.MenuWidth - MenuPadding).
			Height(m.height - 2).
			Render(utils.Ternary(m.mandelbortModel.displayImg, utils.PadEmptyLines(m.mandelbortModel.image, m.params.Height), m.mandelbortModel.text))

		menuPanel := panelStyle.Render(menuContent)

		return lipgloss.JoinHorizontal(
			lipgloss.Center,
			mandelbrotPanel,
			lipgloss.NewStyle().Width(1).Render(""),
			menuPanel,
		)
	} else {
		return utils.Ternary(m.mandelbortModel.displayImg, utils.PadEmptyLines(m.mandelbortModel.image, m.params.Height), m.mandelbortModel.text)
	}
}

func (m Model) UpdateMandelbrot(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		key := msg.String()
		for action, keys := range keyBindings {
			for _, k := range keys {
				if k == key {
					if action == Quit {
						return m, tea.Quit
					}
					if !m.mandelbortModel.displayImg || action == ToggleImg || action == Hide {
						mandelbrotKeyHandlers[action](&m)
					}
					break
				}
			}
		}
	}
	m.RedrawMandelbrot()
	return m, nil
}
