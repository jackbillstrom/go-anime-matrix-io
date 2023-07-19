package gui

import "fyne.io/fyne/v2"

// Screen defines the data structure for a tutorial
type Screen struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

var (
	// Screens defines the metadata for each screen
	Screens = map[string]Screen{
		"preview": {"Preview", "", welcomeScreen, true},
		"settings": {"Settings",
			"Change theme presets and customize options",
			makeSettingsTab,
			true,
		},
		"advanced": {"Advanced",
			"Debug and advanced information.",
			advancedScreen,
			true,
		},
	}

	// ScreenIndex defines how our tutorials should be laid out in the index tree
	ScreenIndex = map[string][]string{
		"": {"preview", "settings", "advanced"},
	}
)
