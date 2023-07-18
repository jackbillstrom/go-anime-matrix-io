package main

import (
	"fyne.io/fyne/v2/theme"
	"go-anime-matrix-io/internal/gui"
	"go-anime-matrix-io/pkg/utils"
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"

	"embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

//go:embed Icon.png
var IconFile embed.FS

// Main window
var _ fyne.Window

const preferenceCurrentTutorial = "currentTutorial"

// makeTray renders the system tray menu
func makeTray(a fyne.App, w fyne.Window) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("Open settings", func() {
			log.Println("System tray menu tapped")
			w.Show()
		})
		h.Icon = theme.SettingsIcon()
		menu := fyne.NewMenu("Hello World", h)
		desk.SetSystemTrayMenu(menu)
	}
}

func main() {
	utils.DisableAnime()

	defer utils.HandleCrash()

	a := app.NewWithID("com.jackbillstrom.go-anime-matrix-io")
	a.SetIcon(utils.AppIcon(IconFile))

	//logLifecycle(a)
	w := a.NewWindow("Go Anime Matrix")

	_ = w
	makeTray(a, w)

	//w.SetMainMenu(makeMenu(a, w))
	w.SetMaster()

	// Set the window close intercept to hide the window instead of exiting the app
	w.SetCloseIntercept(func() {
		w.Hide()
	})

	content := container.NewMax()
	title := widget.NewLabel("")
	intro := widget.NewLabel("")
	intro.Wrapping = fyne.TextWrapWord
	setTutorial := func(t gui.Screen) {
		title.SetText(t.Title)
		intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setTutorial, false))
	} else {
		split := container.NewHSplit(makeNav(setTutorial, true), tutorial)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}

// makeNav creates the navigation tree for the application list menu
func makeNav(setTutorial func(tutorial gui.Screen), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return gui.ScreenIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := gui.ScreenIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := gui.Screens[uid]
			if !ok {
				fyne.LogError("Missing panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			obj.(*widget.Label).TextStyle = fyne.TextStyle{}
		},
		OnSelected: func(uid string) {
			if t, ok := gui.Screens[uid]; ok {
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	return container.NewBorder(nil, nil, nil, nil, tree)
}
