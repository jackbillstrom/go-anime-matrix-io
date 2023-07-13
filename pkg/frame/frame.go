package frame

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"io"
	"io/fs"
	"log"
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
func NewFrame(width, height, fontSize int, fontFile fs.FS) *Frame {
	img := image.NewPaletted(image.Rect(0, 0, width, height), []color.Color{color.Black, color.White})
	face, err := loadFontFace(fontFile, fontSize)
	if err != nil {
		log.Fatalf("Failed to load font: %v", err)
	}

	return &Frame{
		Img:   img,
		Delay: Delay,
		Face:  face,
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

func loadFontFace(fontFile fs.FS, size int) (font.Face, error) {
	f, err := fontFile.Open("static/pixelmix.ttf")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the .ttf file
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	font, err := opentype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: float64(size),
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}

	return face, nil
}
