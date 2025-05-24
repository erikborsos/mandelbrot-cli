package mandelbrot

import (
	"bytes"
	"image"
	"image/png"
	"math"
	"math/cmplx"
	"strings"
	"sync"
)

// MandelbrotParams holds parameters for rendering the Mandelbrot set.
type MandelbrotParams struct {
	CenterRe, CenterIm float64
	ZoomFactor         float64
	MaxIter            int
	Width, Height      int
	ColorMode          int
	Smooth             bool
}

// Reset sets parameters back to default, keeping size intact
func (p *MandelbrotParams) Reset() {
	p.Overwrite(InitialMandelbrotParams())
}

func (p *MandelbrotParams) Overwrite(new MandelbrotParams) {
	p.CenterRe = new.CenterRe
	p.CenterIm = new.CenterIm
	p.ZoomFactor = new.ZoomFactor
	p.MaxIter = new.MaxIter
}

// Move moves the center of the view
func (p *MandelbrotParams) Move(dx, dy float64) {
	p.CenterRe += dx * p.ZoomFactor
	p.CenterIm += dy * p.ZoomFactor
}

// ZoomIn zooms in by reducing zoom factor
func (p *MandelbrotParams) ZoomIn() {
	p.ZoomFactor *= 0.75
}

// ZoomOut zooms out by increasing zoom factor
func (p *MandelbrotParams) ZoomOut() {
	p.ZoomFactor /= 0.74
}

// CycleColor cycles through color modes
func (p *MandelbrotParams) CycleColor() {
	p.ColorMode = (p.ColorMode + 1) % ColorModeCount
}

// ToggleSmooth toggles smooth coloring on/off
func (p *MandelbrotParams) ToggleSmooth() {
	p.Smooth = !p.Smooth
}

// IncreaseIterations adds 10 to max iterations
func (p *MandelbrotParams) IncreaseIterations() {
	p.MaxIter += 10
}

// DecreaseIterations subtracts 10 from max iterations, minimum 10
func (p *MandelbrotParams) DecreaseIterations() {
	if p.MaxIter > 10 {
		p.MaxIter -= 10
	}
}

func InitialMandelbrotParams() MandelbrotParams {
	return MandelbrotParams{
		CenterRe:   -0.5,
		CenterIm:   0,
		ZoomFactor: 1.0,
		MaxIter:    100,
		ColorMode:  ColorNebula,
		Smooth:     true,
	}
}

// mandelbrot computes the number of iterations before divergence for point c.
func mandelbrot(c complex128, maxIter int, smooth bool) float64 {
	z := complex(0, 0)
	for i := range maxIter {
		z = z*z + c
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			if smooth {
				smoothIter := float64(i) - math.Log(math.Log(cmplx.Abs(z)))/math.Log(2)
				smoothNorm := math.Mod(smoothIter/float64(maxIter), 1.0)
				return smoothNorm * float64(maxIter)
			}
			return float64(i)
		}
	}
	return float64(maxIter)
}

// generateMandelbrotText generates the Mandelbrot set as a string buffer.
func GenerateMandelbrotText(params MandelbrotParams) [][]string {
	scale := 3.25 * params.ZoomFactor
	minRe := params.CenterRe - scale/2
	maxRe := params.CenterRe + scale/2
	aspectRatio := float64(params.Height) / float64(params.Width)
	minIm := params.CenterIm - scale*aspectRatio/2
	maxIm := params.CenterIm + scale*aspectRatio/2

	buffer := make([][]string, params.Height)

	var wg sync.WaitGroup
	for y := 0; y < params.Height; y++ {
		y := y // capture loop variable
		wg.Add(1)
		go func() {
			buffer[y] = make([]string, params.Width)
			for x := 0; x < params.Width; x++ {
				c := complex(
					minRe+float64(x)*(maxRe-minRe)/float64(params.Width),
					minIm+float64(y)*(maxIm-minIm)/float64(params.Height),
				)
				iterations := mandelbrot(c, params.MaxIter, params.Smooth)
				buffer[y][x] = getColorString(getColor(params.ColorMode, iterations, params.MaxIter))
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return buffer
}

func BufferToString(buffer [][]string) string {
	var mandelbrotBuilder strings.Builder
	for y := range buffer {
		for x := range buffer[y] {
			mandelbrotBuilder.WriteString(buffer[y][x])
		}
		if y < len(buffer)-1 {
			mandelbrotBuilder.WriteByte('\n')
		}
	}
	return mandelbrotBuilder.String()
}

// generateMandelbrotImage creates a PNG image of Mandelbrot
// width and height can be larger than text buffer, but keep aspect ratio same.
func GenerateMandelbrotImage(params MandelbrotParams) ([]byte, error) {
	aspectRatio := float64(params.Height) / float64(params.Width)
	imgWidth := 1920
	imgHeight := int(float64(imgWidth) * aspectRatio)
	scale := 3.25 * params.ZoomFactor
	minRe := params.CenterRe - scale/2
	maxRe := params.CenterRe + scale/2
	minIm := params.CenterIm - scale*aspectRatio/2
	maxIm := params.CenterIm + scale*aspectRatio/2

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	deltaRe := (maxRe - minRe) / float64(imgWidth)
	deltaIm := (maxIm - minIm) / float64(imgHeight)

	var wg sync.WaitGroup
	for y := range imgHeight {
		y := y // capture loop variable
		wg.Add(1)
		go func() {
			defer wg.Done()
			im := minIm + float64(y)*deltaIm
			for x := range imgWidth {
				re := minRe + float64(x)*deltaRe
				c := complex(re, im)
				iter := mandelbrot(c, params.MaxIter, params.Smooth)
				col := getColor(params.ColorMode, iter, params.MaxIter)
				img.Set(x, y, col)
			}
		}()
	}

	wg.Wait()

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
