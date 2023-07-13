package utils

import (
	"log"
	"os/exec"
)

// EnableAnime enables the anime matrix display via a bash-command
func EnableAnime() {
	cmd := exec.Command("asusctl", "anime", "-e", "true")
	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed to enable anime: ", err)
	}
}

// DisableAnime disables the anime matrix display, used at os.Exit
func DisableAnime() {
	cmd := exec.Command("asusctl", "anime", "-e", "false")
	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed to disable anime: ", err)
	}
}

// Display is used for asking asusctl to use out generated gif-image as output on the display
func Display(fileName string) {
	cmd := exec.Command("asusctl", "anime", "pixel-gif", "-p", fileName)
	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed to display graphics: ", err)
	}
}