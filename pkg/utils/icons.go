package utils

import (
	"bytes"
	"embed"
	"image/png"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/v2"
)

func AppIcon(IconFile embed.FS) fyne.Resource {
	// Load the icon from embed FS
	iconBytes, err := IconFile.ReadFile("Icon.png")
	if err != nil {
		fyne.LogError("Failed to read icon", err)
		return nil // or return some default icon
	}

	iconImage, err := png.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		fyne.LogError("Failed to decode icon", err)
		return nil // or return some default icon
	}

	return canvas.NewImageFromImage(iconImage).Resource
}
