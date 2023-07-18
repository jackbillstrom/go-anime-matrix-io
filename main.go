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

//go:embed static/fonts/pixelmix.ttf
var FontFile embed.FS

//go:embed Icon.png
var IconFile embed.FS

// Main window
var topWindow fyne.Window

const preferenceCurrentTutorial = "currentTutorial"

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
	a := app.NewWithID("github.jackbillstrom.go-anime-matrix-io")
	a.SetIcon(utils.AppIcon(IconFile))

	//logLifecycle(a)
	w := a.NewWindow("Go Anime Matrix")
	topWindow = w

	makeTray(a, w)

	//w.SetMainMenu(makeMenu(a, w))
	w.SetMaster()

	// Set the window close intercept to hide the window instead of exiting the app
	w.SetCloseIntercept(func() {
		w.Hide()
	})

	content := container.NewMax()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	setTutorial := func(t gui.Screen) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

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
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := gui.Screens[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedTutorial(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := gui.Screens[uid]; ok {
				if unsupportedTutorial(t) {
					return
				}
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

func unsupportedTutorial(t gui.Screen) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}
