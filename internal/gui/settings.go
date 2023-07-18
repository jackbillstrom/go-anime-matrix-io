package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AppSettings struct {
	Mode            string
	ShowCPUTemp     bool
	CPUFanSpeed     string
	CPULoadOrRAMUse string
	ShowSongTitle   bool
	ShowEqualizer   bool
	Brightness      float64
}

// makeSettingsTab creates a view for accessing options
func makeSettingsTab(_ fyne.Window) fyne.CanvasObject {
	appSettings := &AppSettings{}
	// Mode select
	modeSelect := widget.NewSelect([]string{"System mode", "Audio mode"}, func(s string) {
		appSettings.Mode = s
		fmt.Println("selected mode", s)
	})

	// CPU settings
	cpuTempCheck := widget.NewCheck("Show CPU temperature", func(on bool) {
		appSettings.ShowCPUTemp = on
		fmt.Println("checked CPU Temp", on)
	})

	// CPU fan speed
	cpuFanSpeedSelect := widget.NewSelect([]string{"CPU Fan Speed", "GPU Fan Speed", "Average Fan Speeds", "Battery"}, func(s string) {
		appSettings.CPUFanSpeed = s
		fmt.Println("selected CPU fan speed", s)
	})

	// CPU Load or RAM usage
	cpuLoadOrRAMUse := widget.NewRadioGroup([]string{"CPU Load", "RAM usage"}, func(s string) {
		appSettings.CPULoadOrRAMUse = s
		fmt.Println("selected CPU load or RAM usage", s)
	})

	// Audio settings
	showSongTitleCheck := widget.NewCheck("Show song title", func(on bool) {
		appSettings.ShowSongTitle = on
		fmt.Println("checked Show Song Title", on)
	})

	showEqualizerCheck := widget.NewCheck("Show equalizer", func(on bool) {
		appSettings.ShowEqualizer = on
		fmt.Println("checked Show Equalizer", on)
	})

	// Brightness slider
	brightnessSlider := widget.NewSlider(0, 100)
	brightnessSlider.Value = 100.0
	brightnessSlider.Step = 10
	brightnessSlider.OnChanged = func(value float64) {
		appSettings.Brightness = value
		fmt.Println("Brightness changed", value)
	}
	// Labels
	themeLabel := widget.NewLabelWithStyle("Select a theme preset", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	fanLabel := widget.NewLabelWithStyle("Select a fan", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	mainContainer := container.NewVBox(
		themeLabel,
		modeSelect,
		fanLabel,
		cpuTempCheck,
		cpuFanSpeedSelect,
		cpuLoadOrRAMUse,
		showSongTitleCheck,
		showEqualizerCheck,
		brightnessSlider,
	)

	modeSelect.OnChanged = func(s string) {
		// Clear the container and add modeSelect
		mainContainer.Objects = []fyne.CanvasObject{modeSelect}

		// Now add the widgets that are relevant for the selected mode
		if s == "System mode" {
			mainContainer.Add(cpuTempCheck)
			mainContainer.Add(fanLabel)
			mainContainer.Add(cpuFanSpeedSelect)
			mainContainer.Add(cpuLoadOrRAMUse)
		} else if s == "Audio mode" {
			mainContainer.Add(showSongTitleCheck)
			mainContainer.Add(showEqualizerCheck)
		}
		mainContainer.Add(brightnessSlider) // Always show brightness slider
		mainContainer.Refresh()             // Refresh the container to show updated widgets
	}

	return mainContainer
}
