package gui

import (
	"context"
	"fyne.io/fyne/v2/theme"
	"go-anime-matrix-io/internal/models"
	"go-anime-matrix-io/pkg/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// makeSettingsTab creates a view for accessing options
func makeSettingsTab(_ fyne.Window) fyne.CanvasObject {
	var preferences = fyne.CurrentApp().Preferences()
	var cancelFunc context.CancelFunc // Saving cancel function for later use
	appSettings := &models.AppSettings{}

	ctx, cancel := context.WithCancel(context.Background())
	cancelFunc = cancel

	playAction := widget.NewToolbarAction(theme.MediaPlayIcon(), func() {})

	// Toolbar
	toolBar := widget.NewToolbar(widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		preferences.SetString("anim_mode", appSettings.Mode)
		preferences.SetBool("show_cpu_temp", appSettings.ShowCPUTemp)
		preferences.SetString("show_fan_speed", appSettings.CPUFanSpeed)
	}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarSpacer(),
		playAction,
		widget.NewToolbarAction(theme.MediaStopIcon(), func() {
			appSettings.Enabled = false
			utils.DisableAnime()
			playAction.Icon = theme.MediaPlayIcon()
		}),
	)

	// Play/pause the animation
	playAction.OnActivated = func() {
		if appSettings.Enabled {
			appSettings.Enabled = false
			cancelFunc()
			utils.DisableAnime() // Disable animation just in case it's still running
			playAction.Icon = theme.MediaPlayIcon()
		} else {
			// Enable
			appSettings.Enabled = true
			go func() {
				cancelFunc, _ = utils.Startup(
					ctx,
					appSettings,
				)
			}()
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
	// By default, show CPU temperature
	cpuTempCheck.Checked = true
	appSettings.ShowCPUTemp = true

	// Fan Speeds
	cpuFanSpeedSelect := widget.NewSelect([]string{"CPU Fan Speed", "GPU Fan Speed", "Average Fan Speeds"}, func(s string) {
		appSettings.CPUFanSpeed = s
	})
	// By default, show CPU Fan Speed
	cpuFanSpeedSelect.SetSelected("CPU Fan Speed")
	appSettings.CPUFanSpeed = "CPU Fan Speed"

	// CPU Load or RAM usage
	cpuLoadOrRAMUse := widget.NewRadioGroup([]string{"CPU Load", "RAM usage"}, func(s string) {
		appSettings.CPULoadOrRAMUse = s
	})
	// By default, show CPU Load
	cpuLoadOrRAMUse.SetSelected("CPU Load")
	appSettings.CPULoadOrRAMUse = "CPU Load"

	// Audio settings
	showSongTitleCheck := widget.NewCheck("Show song title", func(on bool) {
		appSettings.ShowSongTitle = on
	})
	showEqualizerDemoCheck := widget.NewCheck("Show equalizer demo", func(on bool) {
		appSettings.EqualizerDemo = on
	})

	// TODO: Brightness slider
	brightnessSlider := widget.NewSlider(0, 100)
	brightnessSlider.Value = 100.0
	brightnessSlider.Step = 10
	brightnessSlider.OnChanged = func(value float64) {
		appSettings.Brightness = value
	}

	// Labels
	themeLabel := widget.NewLabelWithStyle("Select a theme preset", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	sensorLabel := widget.NewLabelWithStyle("Select sensor", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Create the main container
	settingsLayout := container.NewVBox(
		themeLabel,
		modeSelect,
	)

	// Update visible widgets based on selected mode
	modeSelect.OnChanged = func(s string) {
		appSettings.Mode = s
		// Clear the container and add modeSelect
		settingsLayout.Objects = []fyne.CanvasObject{modeSelect}

		// Now add the widgets that are relevant for the selected mode
		if s == "System mode" {
			settingsLayout.Add(cpuTempCheck)
			settingsLayout.Add(sensorLabel)
			settingsLayout.Add(cpuFanSpeedSelect)
			settingsLayout.Add(cpuLoadOrRAMUse)
		} else if s == "Audio mode" {
			settingsLayout.Add(showEqualizerDemoCheck)
			settingsLayout.Add(showSongTitleCheck)
		}
		settingsLayout.Refresh() // Refresh the container to show updated widgets
	}

	// A layout that contains the toolbar and the settings layout
	fullLayout := container.NewBorder(toolBar, settingsLayout, nil, nil)
	return fullLayout
}
