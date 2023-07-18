package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"embed"

	"fyne.io/fyne/v2/widget"
)

//go:embed static/fonts/pixelmix.ttf
var FontFile embed.FS

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Anime Matrix IO")

	entry := widget.NewEntry()
	textArea := widget.NewMultiLineEntry()

	settingsForm := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Entry", Widget: entry}},
		OnSubmit: func() { // optional, handle form submission
			log.Println("Form submitted:", entry.Text)
			log.Println("multiline:", textArea.Text)
			myWindow.Close()
		},
	}

	// we can also append items
	settingsForm.Append("Text", textArea)

	// Set up tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Preview", widget.NewLabel("Skall visa GIF-bilden animerad i realtid?")),
		container.NewTabItem("Settings", settingsForm),
		container.NewTabItem("About", widget.NewLabel("TODO")),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
