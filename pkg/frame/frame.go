package frame

import (
	"image"
	"image/color"
	_ "image/draw"

	_ "golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"io/ioutil"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Constants
const (
	FrameWidth  = 64
	FrameHeight = 32
	FrameDelay  = 5  // Lower this to make the animation faster
	FontSize = 11
)

// Frame represents a single frame in an animation
type Frame struct {
	Img   *image.Paletted
	Delay int
	Face  font.Face
	FontSize int
}

// NewFrame initializes a new frame with given dimensions

func NewFrame(width, height int, fontPath string, fontSize float64) *Frame {
	img := image.NewPaletted(image.Rect(0, 0, width, height), []color.Color{color.Black, color.White})
	return &Frame{
		Img:   img,
		Delay: FrameDelay,
		Face:  LoadFontFace(fontPath, fontSize),
		}
}

// DrawText draws the given text onto the frame at the given offset
func (f *Frame) DrawText(text string, row int) {
	// Calculate the y-position based on the row number and the font size
	yPos := (row + 1) * FontSize
	point := fixed.Point26_6{fixed.Int26_6(0 * FrameWidth), fixed.Int26_6(yPos * FrameWidth)}
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
		point := fixed.Point26_6{fixed.Int26_6(offsets[i] * FrameWidth), fixed.Int26_6((i+1) * 13 * FrameWidth)}
		d := &font.Drawer{
			Dst:  f.Img,
			Src:  image.White,
			Face: f.Face,
			Dot:  point,
		}
		d.DrawString(text)
	}
}


func LoadFontFace(path string, size float64) font.Face {
	fontBytes, err := ioutil.ReadFile(path)
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
