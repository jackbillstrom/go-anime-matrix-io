package textscroller

import (
	"go-anime-matrix-io/pkg/frame"
)

// ScrollText generates frames for a scrolling text animation
func ScrollText(text string, numFrames int) []*frame.Frame {
	frames := make([]*frame.Frame, numFrames)
	for i := 0; i < numFrames; i++ {
		f := frame.NewFrame(frame.FrameWidth, frame.FrameHeight)
		f.DrawText(text, i-numFrames/2)
		frames[i] = f
	}
	return frames
}
