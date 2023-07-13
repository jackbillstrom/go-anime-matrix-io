package main

import (
	"go-anime-matrix-io/pkg/gifcreator"
	"go-anime-matrix-io/pkg/textscroller"
)

// Constants
const (
	frameWidth  = 64
)

func main() {
	frames := textscroller.ScrollText("31 / 70 °C", frameWidth)
	frames = append(frames, textscroller.ScrollText("2900 RPM", frameWidth)...)
	frames = append(frames, textscroller.ScrollText("702 KB/s", frameWidth)...)
	err := gifcreator.SaveGif("out.gif", frames)
	if err != nil {
		panic(err)
	}
}