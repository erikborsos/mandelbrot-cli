package tui

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"mandel-cli/mandelbrot"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Resolution defines a resolution option with a name and dimensions.
type Resolution struct {
	Name   string
	Width  int
	Height int
}

// resolutionOptions defines predefined resolutions for the select field.
var resolutionOptions = []Resolution{
	{Name: "720p", Width: 1280, Height: 720},
	{Name: "1080p", Width: 1920, Height: 1080},
	{Name: "WQHD (1440p)", Width: 2560, Height: 1440},
	{Name: "4K", Width: 3840, Height: 2160},
}

type SaveModel struct {
	form      *huh.Form
	errorMsg  string
	completed bool
}

func initSaveModel(params mandelbrot.MandelbrotParams) SaveModel {
	// Initialize default resolution (match params.Width and params.Height or default to 1080p)
	defaultResolution := "1080p"
	for _, res := range resolutionOptions {
		if res.Width/2 == params.Width && res.Height == params.Height {
			defaultResolution = res.Name
			break
		}
	}

	// Initialize color options
	var colorOptions []huh.Option[string]
	for _, color := range mandelbrot.ColorNames {
		colorOptions = append(colorOptions, huh.NewOption(color, color))
	}

	// Initialize resolution options for select field
	var resOptions []huh.Option[string]
	for _, res := range resolutionOptions {
		resOptions = append(resOptions, huh.NewOption(res.Name, res.Name))
	}

	// Initialize file path and color
	filepathStr := "mandelbrot.png"
	colorStr := mandelbrot.ColorNames[params.ColorMode]

	// Create form with resolution select
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Resolution").
				Key("resolution").
				Options(resOptions...).
				Value(&defaultResolution),
			huh.NewSelect[string]().
				Title("Color Scheme").
				Key("color").
				Options(colorOptions...).
				Value(&colorStr),
			huh.NewInput().
				Title("File Path").
				Key("filepath").
				Value(&filepathStr).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("file path cannot be empty")
					}
					if !strings.HasSuffix(strings.ToLower(s), ".png") {
						return fmt.Errorf("file path must end with .png")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm())
	form.Init()

	return SaveModel{
		form:      form,
		errorMsg:  "",
		completed: false,
	}
}

func (m Model) UpdateSave(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.view = MandelbrotView
			m.saveModel = initSaveModel(m.params) // Reset form on exit
			return m, nil
		}

	case tea.WindowSizeMsg:
		// Update form size
		h, v := docStyle.GetFrameSize()
		m.saveModel.form.WithWidth(msg.Width - h).WithHeight(msg.Height - v)
	}

	var cmd tea.Cmd
	var model tea.Model
	model, cmd = m.saveModel.form.Update(msg)
	if form, ok := model.(*huh.Form); ok {
		m.saveModel.form = form
	} else {
		m.saveModel.errorMsg = "Error: Failed to update form"
		return m, nil
	}

	if m.saveModel.form.State == huh.StateCompleted {
		if !m.saveModel.completed {
			m.saveModel.completed = true
			resolution := m.saveModel.form.GetString("resolution")
			filepath := m.saveModel.form.GetString("filepath")
			color := m.saveModel.form.GetString("color")

			// Find selected resolution's width and height
			var width, height int
			for _, res := range resolutionOptions {
				if res.Name == resolution {
					width = res.Width
					height = res.Height
					break
				}
			}
			if width == 0 || height == 0 {
				m.saveModel.errorMsg = "Invalid resolution selected"
				m.saveModel.completed = false
				return m, nil
			}

			// Update params with new values
			saveParams := m.params
			saveParams.Width = width
			saveParams.Height = height
			for i, name := range mandelbrot.ColorNames {
				if name == color {
					saveParams.ColorMode = i
					break
				}
			}

			// Generate and save the image
			img, err := mandelbrot.GenerateFixedMandelbrotImage(saveParams, saveParams.Width, saveParams.Height)
			if err != nil {
				m.saveModel.errorMsg = fmt.Sprintf("Error generating image: %v", err)
				m.saveModel.completed = false
				return m, nil
			}

			err = SaveImage(img, filepath)
			if err != nil {
				m.saveModel.errorMsg = fmt.Sprintf("Error saving image: %v", err)
				m.saveModel.completed = false
				return m, nil
			}

			// Switch back to Mandelbrot view after successful save
			m.view = MandelbrotView
			m.saveModel = initSaveModel(m.params) // Reset form
			m.saveModel.errorMsg = ""
		}
	}

	return m, cmd
}

func (m Model) ViewSave() string {
	var b strings.Builder
	b.WriteString(docStyle.Render(m.saveModel.form.View()))
	if m.saveModel.errorMsg != "" {
		b.WriteString("\n" + errorStyle.Render("Error: "+m.saveModel.errorMsg))
	}
	return b.String()
}

func SaveImage(imgByte []byte, filepath string) error {
	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalln(err)
	}
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
