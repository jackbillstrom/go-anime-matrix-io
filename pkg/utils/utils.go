package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
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
		DisableAnime()
	}
}

// CheckCommands is used for checking for necessary tools are installed
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

// SetupService is used for setting up the systemd service file
func SetupService() {
	// Ask the user if they want to install the service
	fmt.Print("Do you want to install the service? [Y/n] ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))

	// If the user wants to install the service, create the systemd service file
	if response == "y" || response == "yes" {
		// Get the current user
		currentUser, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		// Get the current working directory
		workingDirectory, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		// Get the current binary path
		binaryPath, err := filepath.Abs(os.Args[0])
		if err != nil {
			log.Fatal(err)
		}

		service := fmt.Sprintf(`[Unit]
Description=go-anime-matrix-io service for using the animated matrix display in the background
After=network.target

[Service]
ExecStart=%s
WorkingDirectory=%s
User=%s
Group=%s
Restart=always

[Install]
WantedBy=multi-user.target
`, binaryPath, workingDirectory, currentUser.Username, currentUser.Username) // Assuming a user's group is the same as their username

		err = os.WriteFile("/etc/systemd/system/anime-matrix-io.service", []byte(service), 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Service installed. You can now enable it with 'systemctl enable anime-matrix-io'.")
	}
}
