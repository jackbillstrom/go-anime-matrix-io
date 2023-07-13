package gifcreator

import (
	"image/gif"
	"os"

	"go-anime-matrix-io/pkg/frame"
)

// SaveGif saves a sequence of frames as a gif to a file
func SaveGif(filename string, frames []*frame.Frame) error {
	outGif := &gif.GIF{}
	for _, f := range frames {
		outGif.Image = append(outGif.Image, f.Img)
		outGif.Delay = append(outGif.Delay, f.Delay)
	}
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return gif.EncodeAll(outFile, outGif)
}
