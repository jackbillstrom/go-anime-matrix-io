package gifcreator

import (
	"image"
	"image/draw"
	"image/gif"
	"os"

	"go-anime-matrix-io/pkg/frame"
)

// SaveGif saves a sequence of frames as a gif to a file
func SaveGif(filename string, frames []*frame.Frame, isMirrored bool) error {
	outGif := &gif.GIF{}
	for _, f := range frames {
		// If isMirrored is true, mirror the image vertically
		if isMirrored {
			flipVertical(f.Img)
		}

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

// flipVertical flips the image vertically
func flipVertical(img image.Image) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Create a new image with the same bounds
	mirroredImg := image.NewRGBA(bounds)

	// Draw the original image onto the new image, flipping it vertically
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			mirroredImg.Set(x, height-y-1, img.At(x, y))
		}
	}

	// Replace the original image with the mirrored image
	draw.Draw(img.(draw.Image), bounds, mirroredImg, image.Point{}, draw.Src)
}
