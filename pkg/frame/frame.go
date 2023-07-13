package frame

import (
	"image"
	"image/color"
	_ "image/draw"
	"os"

	_ "golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Constants
const (
	Width    = 64
	Height   = 32
	Delay    = 5 // Lower this to make the animation faster
	FontSize = 7
)

// Frame represents a single frame in an animation
type Frame struct {
	Img      *image.Paletted
	Delay    int
	Face     font.Face
	FontSize int
}

// NewFrame initializes a new frame with given dimensions
func NewFrame(width, height int, fontPath string, fontSize float64) *Frame {
	img := image.NewPaletted(image.Rect(0, 0, width, height), []color.Color{color.Black, color.White})
	return &Frame{
		Img:   img,
		Delay: Delay,
		Face:  loadFontFace(fontPath, fontSize),
	}
}

// DrawText draws the given text onto the frame at the given offset
func (f *Frame) DrawText(text string, row int) {
	// Calculate the y-position based on the row number and the font size
	yPos := (row + 1) * FontSize
	point := fixed.Point26_6{X: fixed.Int26_6(0 * Width), Y: fixed.Int26_6(yPos * Width)}
	d := &font.Drawer{
		Dst:  f.Img,
		Src:  image.White,
		Face: f.Face,
		Dot:  point,
	}
	d.DrawString(text)
}

// DrawTextMulti makes use of multiple lines of text
func (f *Frame) DrawTextMulti(texts []string, offsets []int) {
	for i, text := range texts {
		point := fixed.Point26_6{X: fixed.Int26_6(offsets[i] * Width), Y: fixed.Int26_6((i + 1) * 13 * Width)}
		d := &font.Drawer{
			Dst:  f.Img,
			Src:  image.White,
			Face: f.Face,
			Dot:  point,
		}
		d.DrawString(text)
	}
}

// DrawProgressBar draws a progress bar with the given progress
func (f *Frame) DrawProgressBar(progress int, row int) {
	bar := ""
	maxProgress := 12
	for i := 0; i < maxProgress; i++ {
		if i < progress {
			bar += "-"
		} else {
			bar += " "
		}
	}
	f.DrawText(bar, row)
}

// loadFontFace loads a font face from a given path
func loadFontFace(path string, size float64) font.Face {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	f, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return face
}
