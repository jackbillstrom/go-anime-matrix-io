package animated

import (
	"go-anime-matrix-io/pkg/frame"
)

// AnimationHandler is the type of our animation functions
type AnimationHandler func(texts []string, numFrames int, frameNum int) *frame.Frame

// ScrollText generates frames for a scrolling text animation
func ScrollText(text string, numFrames int, frameNum int) *frame.Frame {
	f := frame.NewFrame(frame.Width, frame.Height)
	f.DrawText(text, frameNum-numFrames/2)
	return f
}

// BlinkText generates frames for a blinking text animation
func BlinkText(text string, numFrames int, frameNum int) *frame.Frame {
	f := frame.NewFrame(frame.Width, frame.Height)
	if frameNum%2 == 0 {
		f.DrawText(text, 0)
	}
	return f
}

// GenerateFrames generates the frames for the animation using the given handler
func GenerateFrames(handler AnimationHandler, texts []string, numFrames int) []*frame.Frame {
	frames := make([]*frame.Frame, numFrames)
	for i := 0; i < numFrames; i++ {
		frames[i] = handler(texts, numFrames, i)
	}
	return frames
}

// SlideText is for sliding text while using multiple rows of text
func SlideText(texts []string, numFrames int, frameNum int) *frame.Frame {
	x := frame.NewFrame(frame.Width, frame.Height)
	// Create offsets for CPU temp and fan speed text, so they can be shown in the same row
	offsets := make([]int, len(texts))
	for i := range offsets {
		offsets[i] = frameNum - numFrames/2
	}
	x.DrawTextMulti(texts, offsets)
	return x
}

// GenerateSingleRowFrames generates the frames for the animation using the given handler for single row
func GenerateSingleRowFrames(handler AnimationHandler, text1 string, text2 string, numFrames int) []*frame.Frame {
	frames := make([]*frame.Frame, numFrames)
	for i := 0; i < numFrames; i++ {
		frames[i] = handler([]string{text1, text2}, numFrames, i)
	}
	return frames
}

// StaticText generates frames for a static text
func StaticText(texts []string, numFrames int, frameNum int) *frame.Frame {
	x := frame.NewFrame(frame.Width, frame.Height)
	offsets := make([]int, len(texts))

	// Keep the text in the same position for all frames
	for i := range offsets {
		offsets[i] = 0
	}

	x.DrawTextMulti(texts, offsets)
	return x
}
