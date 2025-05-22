package tui

import (
	"mandel-cli/mandelbrot"
	"mandel-cli/utils"
)

type MandelbrotModel struct {
	params        mandelbrot.MandelbrotParams // Parameters for Mandelbrot rendering
	text          string                      // Text representation of the Mandelbrot set
	image         string                      // Kitty terminal image representation
	width         int                         // Terminal width
	height        int                         // Terminal height
	displayImg    bool                        // Whether to display image or text
	paramsChanged bool                        // Whether parameters have changed
	errorMsg      string                      // Error message for UI display
}

var infoReplacer utils.ChainReplacer
var controls string
var controlsDisabled string
