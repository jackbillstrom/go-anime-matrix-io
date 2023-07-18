package gui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func welcomeScreen(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromFile("static/icons/matrix-light.png")
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(256, 256))
	}

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Anime Matrix IO", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		widget.NewLabelWithStyle("Disabled", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabel(""), // whitespace
		container.NewHBox(
			widget.NewHyperlink("github", parseURL("https://github.com/jackbillstrom/go-anime-matrix-io/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("buy me a coffee (or beer)", parseURL("https://bmc.link/jackbillstrom")),
		),
		widget.NewLabel(""), // whitespace
	))
}
