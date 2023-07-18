package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// makeSettingsTab creates a view for accessing options
func makeSettingsTab(_ fyne.Window) fyne.CanvasObject {
	disabledCheck := widget.NewCheck("Disabled check", func(bool) {})
	disabledCheck.Disable()
	checkGroup := widget.NewCheckGroup([]string{"CheckGroup Item 1", "CheckGroup Item 2"}, func(s []string) { fmt.Println("selected", s) })
	checkGroup.Horizontal = true
	radio := widget.NewRadioGroup([]string{"Radio Item 1", "Radio Item 2"}, func(s string) { fmt.Println("selected", s) })
	radio.Horizontal = true
	disabledRadio := widget.NewRadioGroup([]string{"Disabled radio"}, func(string) {})
	disabledRadio.Disable()

	// Brightness slider
	brightnessSliderLabel := widget.NewLabel("Brightness") // TODO: Show current value in realtime
	brightnessSlider := widget.NewSlider(0, 100)
	brightnessSlider.Value = 100.0
	brightnessSlider.Step = 10

	return container.NewVBox(
		widget.NewSelect([]string{"System mode", "Audio mode"}, func(s string) { fmt.Println("selected", s) }),
		widget.NewCheck("Show CPU temperature", func(on bool) { fmt.Println("checked", on) }),
		disabledCheck,
		checkGroup,
		radio,
		disabledRadio,
		brightnessSliderLabel,
		brightnessSlider,
	)
}
