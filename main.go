package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"go-anime-matrix-io/pkg/frame"
	"go-anime-matrix-io/pkg/gifcreator"
	"go-anime-matrix-io/pkg/sensors"
	"go-anime-matrix-io/pkg/utils"
	"time"
)

// Constants for necessary parameters
const (
	frameWidth  = 64 // Size of Anime-matrix display
	numFrames   = 30 // Number of frames for the animations
	fontPath = "./static/Hack-Regular.ttf" // Change this to the actual path of your font
	fontSize = 10
	fileName = "out.gif"
	seconds = 2
)

// handleCrash is run to disable anime matrix after a recovery
func handleCrash() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from crash, disabling anime:", r)
		utils.DisableAnime()
	}
}

func main() {
	// Check for necessary stuff are installed or whatnot
	err := utils.CheckCommands()
	if err != nil {
		fmt.Println("Required command is missing:", err)
		return
	}

	// Clear & enable matrix display via asusctl
	utils.EnableAnime()

	// Create a ticker that triggers every X seconds
	ticker := time.NewTicker(seconds * time.Second)

	// Make sure anime is disabled when the program ends
	defer utils.DisableAnime()

	// Make sure anime is disabled when the program crashes
	defer handleCrash()

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// setup signal catching
	sigs := make(chan os.Signal, 1)

	// register the sigs channel to receive SIGINT
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Launch a goroutine that will update the data and generate the GIF
	go func() {
		for {
			select {
			case <-ctx.Done():
				// If the context is cancelled, return from the goroutine
				return
			case <-ticker.C:
				// Fetching data to input into image
				cpuTemp, _, cpuFan, _, err := sensors.GetSensorData()
				if err != nil {
					fmt.Errorf("error fetching sensors data, have you installed 'lm-sensors' and has it reachable via the command 'sensors'?", err)
					continue // Skip this cycle if there was an error fetching the data
				}

				// Prepare information to be displayed on a single row
				cpuInfo := cpuTemp

				// Generate frames for single row text
				frames := make([]*frame.Frame, numFrames)
				for i := 0; i < numFrames; i++ {
					f := frame.NewFrame(frameWidth, frame.FrameHeight, fontPath, fontSize)
					f.DrawText("   " + cpuInfo, 1)
					f.DrawText(" " + cpuFan, 2)
					frames[i] = f
				}

				// Append output to a img
				err = gifcreator.SaveGif(fileName, frames)
				if err != nil {
					fmt.Errorf("error occured while saving the output gif", err)
					continue // Skip this cycle if there was an error saving the GIF
				}

				// Append to anime matrix
				utils.Display(fileName)
			}
		}
	}()

	// Wait for an interrupt signal
	<-sigs

	// When the signal is received, cancel the context
	cancel()

	// Wait a bit to allow the goroutines to clean up
	time.Sleep(time.Second)
}