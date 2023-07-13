package animated

import (
	"go-anime-matrix-io/pkg/frame"
)

// AnimationHandler is the type of our animation functions
type AnimationHandler func(texts []string, numFrames int, frameNum int) *frame.Frame
