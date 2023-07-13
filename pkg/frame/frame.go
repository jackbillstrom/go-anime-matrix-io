package frame

import (
	"image"
	"image/color"
	_ "image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Constants
const (
	FrameWidth  = 64
	FrameHeight = 32
	FrameDelay  = 5  // Lower this to make the animation faster
)

// Frame represents a single frame in an animation
type Frame struct {
	Img   *image.Paletted
	Delay int
}

// NewFrame initializes a new frame with given dimensions
func NewFrame(width, height int) *Frame {
	img := image.NewPaletted(image.Rect(0, 0, width, height), []color.Color{color.Black, color.White})
	return &Frame{
		Img:   img,
		Delay: FrameDelay,
	}
}

// DrawText draws the given text onto the frame at the given offset
func (f *Frame) DrawText(text string, offset int) {
	point := fixed.Point26_6{fixed.Int26_6(offset * 64), fixed.Int26_6((FrameHeight/2 + 5) * 64)}
	d := &font.Drawer{
		Dst:  f.Img,
		Src:  image.White,
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}