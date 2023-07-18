package gui

import (
	"fmt"
	"fyne.io/fyne/v2/theme"
	"go-anime-matrix-io/pkg/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AppSettings struct {
	Enabled         bool    `json:"enabled"`
	Mode            string  `json:"mode"`
	ShowCPUTemp     bool    `json:"show_cpu_temp"`
	CPUFanSpeed     string  `json:"cpu_fan_speed"`
	CPULoadOrRAMUse string  `json:"cpu_load_or_ram_use"`
	ShowSongTitle   bool    `json:"show_song_title"`
	ShowEqualizer   bool    `json:"show_equalizer"`
	Brightness      float64 `json:"brightness"`
}

// makeSettingsTab creates a view for accessing options
func makeSettingsTab(_ fyne.Window) fyne.CanvasObject {
	var preferences = fyne.CurrentApp().Preferences()
	appSettings := &AppSettings{}

	playAction := widget.NewToolbarAction(theme.MediaPlayIcon(), func() {})

	// Toolbar
	toolBar := widget.NewToolbar(widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		preferences.SetString("anim_mode", appSettings.Mode)
		preferences.SetBool("show_cpu_temp", appSettings.ShowCPUTemp)
		preferences.SetString("show_fan_speed", appSettings.CPUFanSpeed)
		fmt.Printf("%+v\n", appSettings)
	}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarSpacer(),
		playAction,
		widget.NewToolbarAction(theme.MediaReplayIcon(), func() { fmt.Println("Restarting ...") }),
	)

	playAction.OnActivated = func() {
		if appSettings.Enabled {
			go utils.DisableAnime()
			appSettings.Enabled = false
			playAction.Icon = theme.MediaPlayIcon()
		} else {
			// Enable
			go utils.Startup()
			appSettings.Enabled = true
			playAction.Icon = theme.MediaPauseIcon()
		}
		toolBar.Refresh()
	}

	// Mode select
	modeSelect := widget.NewSelect([]string{"System mode", "Audio mode"}, func(s string) {
		appSettings.Mode = s
	})

	// CPU settings
	cpuTempCheck := widget.NewCheck("Show CPU temperature", func(on bool) {
		appSettings.ShowCPUTemp = on
	})

	// Fan Speeds
	cpuFanSpeedSelect := widget.NewSelect([]string{"CPU Fan Speed", "GPU Fan Speed", "Average Fan Speeds", "Battery"}, func(s string) {
		appSettings.CPUFanSpeed = s
	})

	// CPU Load or RAM usage
	cpuLoadOrRAMUse := widget.NewRadioGroup([]string{"CPU Load", "RAM usage"}, func(s string) {
		appSettings.CPULoadOrRAMUse = s
	})

	// Audio settings
	showSongTitleCheck := widget.NewCheck("Show song title", func(on bool) {
		appSettings.ShowSongTitle = on
	})

	// Equalizer
	showEqualizerCheck := widget.NewCheck("Show equalizer", func(on bool) {
		appSettings.ShowEqualizer = on
	})

	// Brightness slider
	brightnessSlider := widget.NewSlider(0, 100)
	brightnessSlider.Value = 100.0
	brightnessSlider.Step = 10
	brightnessSlider.OnChanged = func(value float64) {
		appSettings.Brightness = value
	}

	// Labels
	themeLabel := widget.NewLabelWithStyle("Select a theme preset", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	fanLabel := widget.NewLabelWithStyle("Select sensor", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Create the main container
	settingsLayout := container.NewVBox(
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

	// Update visible widgets based on selected mode
	modeSelect.OnChanged = func(s string) {
		// Clear the container and add modeSelect
		settingsLayout.Objects = []fyne.CanvasObject{modeSelect}

		// Now add the widgets that are relevant for the selected mode
		if s == "System mode" {
			settingsLayout.Add(cpuTempCheck)
			settingsLayout.Add(fanLabel)
			settingsLayout.Add(cpuFanSpeedSelect)
			settingsLayout.Add(cpuLoadOrRAMUse)
		} else if s == "Audio mode" {
			settingsLayout.Add(showSongTitleCheck)
			settingsLayout.Add(showEqualizerCheck)
		}
		settingsLayout.Add(brightnessSlider) // Always show brightness slider
		settingsLayout.Refresh()             // Refresh the container to show updated widgets
	}

	// A layout that contains the toolbar and the settings layout
	fullLayout := container.NewBorder(toolBar, nil, nil, settingsLayout)
	return fullLayout
}
