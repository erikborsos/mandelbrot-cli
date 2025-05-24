package tui

import (
	"mandel-cli/mandelbrot"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

type View int

const (
	MandelbrotView View = iota
	PresetsView
	SaveView
)

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
	ForceQuit    KeyAction = "force_quit"
	Hide         KeyAction = "hide"
	SelectPreset KeyAction = "select_preset"
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
}

type Model struct {
	params          mandelbrot.MandelbrotParams // Parameters for Mandelbrot rendering
	width           int                         // Terminal width
	height          int                         // Terminal heighti
	mandelbortModel MandelbrotModel
	presetsModel    PresetsModel
	view            View
}

func InitModel() Model {
	return Model{
		params:          mandelbrot.InitialMandelbrotParams(),
		mandelbortModel: initMandelbrotModel(),
		presetsModel:    initPresetsModel(),
		view:            MandelbrotView,
	}
}

func (m Model) Init() tea.Cmd {
	setupUIConfig()
	return nil
}

func (m *Model) Reset() {
	currentWidth := m.params.Width
	currentHeight := m.params.Height
	m.params.Reset()
	m.params.Width = currentWidth
	m.params.Height = currentHeight
	m.mandelbortModel.errorMsg = ""
}

func (m *Model) RedrawMandelbrot() {
	if !m.mandelbortModel.displayImg && m.mandelbortModel.paramsChanged {
		m.mandelbortModel.text = mandelbrot.BufferToString(mandelbrot.GenerateMandelbrotText(m.params))
		m.mandelbortModel.paramsChanged = false
	} else if m.mandelbortModel.displayImg && m.mandelbortModel.paramsChanged {
		m.toggleDisplayImg()
		m.toggleDisplayImg()
		m.mandelbortModel.paramsChanged = false

	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if slices.Contains(keyBindings[ForceQuit], msg.String()) {
			return m, tea.Quit
		}
	}

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = msg.Width
		m.height = msg.Height
		if m.mandelbortModel.hideMenu {
			m.params.Width = m.width / 2
		} else {
			m.params.Width = (m.width-defaultUIConfig.MenuWidth)/2 - WidthAdjustment
		}
		m.params.Height = m.height
		m.mandelbortModel.paramsChanged = true
	}

	if m.view == MandelbrotView {
		return m.UpdateMandelbrot(msg)
	} else if m.view == PresetsView {
		return m.UpdatePresets(msg)
	}
	return m, nil
}

func (m Model) View() string {
	if m.view == MandelbrotView {
		return m.ViewMandelbrot()
	} else if m.view == PresetsView {
		return m.ViewPresets()
	}
	return ""
}
