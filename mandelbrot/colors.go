package mandelbrot

import (
	"fmt"
	"image/color"
	"math"
)

// Constants for color modes
const (
	ColorGrayscale = iota
	ColorNebula
	ColorRainbow
	ColorFire
	ColorOcean
	ColorPsychedelic
	ColorIce
	ColorInferno
	ColorDesert
	ColorForest
	ColorModeCount
)

var ColorNames = map[int]string{
	ColorGrayscale:   "Grayscale",
	ColorRainbow:     "Rainbow",
	ColorFire:        "Fire",
	ColorPsychedelic: "Psychedelic",
	ColorOcean:       "Ocean",
	ColorIce:         "Ice",
	ColorInferno:     "Inferno",
	ColorDesert:      "Desert",
	ColorNebula:      "Nebula",
	ColorForest:      "Forest",
}

// hsvToRGBA converts HSV values (h in [0,360], s,v in [0,1]) to color.RGBA
func hsvToRGBA(h float64, s float64, v float64) color.Color {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c

	var r, g, b float64

	switch {
	case h >= 0 && h < 60:
		r, g, b = c, x, 0
	case h >= 60 && h < 120:
		r, g, b = x, c, 0
	case h >= 120 && h < 180:
		r, g, b = 0, c, x
	case h >= 180 && h < 240:
		r, g, b = 0, x, c
	case h >= 240 && h < 300:
		r, g, b = x, 0, c
	case h >= 300 && h < 360:
		r, g, b = c, 0, x
	default:
		r, g, b = 0, 0, 0
	}

	r = (r + m) * 255
	g = (g + m) * 255
	b = (b + m) * 255

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

// getColorString returns a string that colors a 2-space block using 24-bit RGB ANSI escape codes
func getColorString(color color.Color) string {
	r, g, b, _ := color.RGBA()
	return fmt.Sprintf("\033[48;2;%d;%d;%dm  \033[0m", r, g, b)
}

// getColor returns a color.Color for a given iteration count and color scheme.
func getColor(scheme int, iterations float64, maxIter int) color.Color {
	if int(iterations) == maxIter {
		return color.RGBA{0, 0, 0, 255}
	}

	// Normalize t to avoid extreme values
	t := math.Max(0, math.Min(1.0, iterations/float64(maxIter)))

	switch scheme {
	case ColorRainbow: // Rainbow
		hue := 360.0 * t
		return hsvToRGBA(hue, 1.0, 1.0)
	case ColorFire: // Fire
		r := math.Min(1.0, t*3)
		g := math.Max(0, math.Min(1.0, t*3-1))
		b := math.Max(0, math.Min(1.0, t*3-2))
		return color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
	case ColorOcean: // Ocean
		r := math.Max(0, math.Min(1.0, t*3-2))
		g := math.Max(0, math.Min(1.0, t*2-1))
		b := math.Min(1.0, t+0.2)
		return color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
	case ColorPsychedelic:
		hue := math.Mod(t*360*5, 360) // cycle hue rapidly
		return hsvToRGBA(hue, 1.0, 1.0)
	case ColorIce:
		r := uint8(0)
		g := uint8(t * 200)
		b := uint8(100 + t*155)
		return color.RGBA{r, g, b, 255}
	case ColorInferno:
		r := uint8(math.Min(1.0, t*4) * 255)
		g := uint8(math.Pow(t, 1.5) * 100)
		b := uint8(math.Pow(1-t, 3) * 255)
		return color.RGBA{r, g, b, 255}
	case ColorDesert:
		r := uint8(255 * t)
		g := uint8(200 * (1 - t))
		b := uint8(100 + 100*t)
		return color.RGBA{r, g, b, 255}
	case ColorNebula:
		hue := 240 + 120*math.Sin(t*4*math.Pi)
		return hsvToRGBA(hue, 0.6+0.4*t, 0.8)
	case ColorForest:
		r := uint8(30 + 50*(1-t))
		g := uint8(100 + 155*t)
		b := uint8(30 + 20*t)
		return color.RGBA{r, g, b, 255}

	default: // Grayscale
		brightness := 1.0 - t
		gray := uint8(brightness * 255)
		return color.RGBA{gray, gray, gray, 255}
	}
}
