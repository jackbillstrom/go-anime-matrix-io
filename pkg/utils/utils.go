package utils

import (
	"fmt"
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

// CheckCommands is used for checking for necessesary tools are installed
func CheckCommands() error {
	commands := []string{
		"asusctl",
		"sensors",
		}

		for _, command := range commands {
			_, err := exec.LookPath(command)
			if err != nil {
				return fmt.Errorf("command not found: %s", command)
			}
		}

		return nil
}