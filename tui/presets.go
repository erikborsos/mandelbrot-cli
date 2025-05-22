package tui

import (
	"mandel-cli/mandelbrot"
)

var Presets = map[string]mandelbrot.MandelbrotParams{
	"Julia Island": {
		CenterRe:   -1.768778770,
		CenterIm:   -0.001738942,
		ZoomFactor: 0.000000340,
		MaxIter:    400,
	},
}
