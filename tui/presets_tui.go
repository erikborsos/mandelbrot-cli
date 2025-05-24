package tui

import (
	"fmt"
	"mandel-cli/mandelbrot"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var presets = map[string]mandelbrot.MandelbrotParams{
	"Julia Island": {
		CenterRe:   -1.768778770,
		CenterIm:   -0.001738942,
		ZoomFactor: 0.000000340,
		MaxIter:    400,
	},
	"Seahorse Valley": {
		CenterRe:   -0.743517833,
		CenterIm:   -0.127094578,
		ZoomFactor: 0.004228283,
		MaxIter:    400,
	},
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// Key map for list actions
type presetListKeyMap struct {
	Select key.Binding
}

var presetKeys = presetListKeyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type PresetsModel struct {
	list list.Model
}

func initPresetsModel() PresetsModel {
	items := make([]list.Item, 0, len(presets))
	for preset := range presets {
		items = append(items, item{
			title: preset,
			desc:  fmt.Sprintf("Real: %v, Imaginary: %v", presets[preset].CenterRe, presets[preset].CenterIm),
		})
	}

	delegate := list.NewDefaultDelegate()
	presetList := list.New(items, delegate, 0, 0)
	presetList.DisableQuitKeybindings()

	// Add custom help keys
	presetList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{presetKeys.Select}
	}
	presetList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{presetKeys.Select}
	}
	presetList.Title = "Select Preset:"
	return PresetsModel{list: presetList}
}

func (m Model) UpdatePresets(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.presetsModel.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if selected, ok := m.presetsModel.list.SelectedItem().(item); ok {
				m.params.Overwrite(presets[selected.title])
				m.view = MandelbrotView
				m.mandelbortModel.paramsChanged = true
			}
		case "esc":
			m.view = MandelbrotView
		}
	}

	var cmd tea.Cmd
	m.presetsModel.list, cmd = m.presetsModel.list.Update(msg)
	m.RedrawMandelbrot()
	return m, cmd
}

func (m Model) ViewPresets() string {
	return docStyle.Render(m.presetsModel.list.View())
}
