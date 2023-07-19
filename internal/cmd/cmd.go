package cmd

import (
	"context"
	"embed"
	"fmt"
	"go-anime-matrix-io/pkg/frame"
	"go-anime-matrix-io/pkg/gifcreator"
	"go-anime-matrix-io/pkg/sensors"
	"go-anime-matrix-io/pkg/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getlantern/systray"
)

// Constants for necessary parameters
const (
	numFrames = 30 // Number of frames for the animations
	fontSize  = 7
	fileName  = "out.gif"
	seconds   = 2
)

func OnReady(fontFile embed.FS) {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Set up signal catching
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to receive quit events from systray
	quitCh := make(chan struct{})

	// Launch a goroutine that will update the data and generate the GIF
	go func() {
		for {
			select {
			case <-ctx.Done():
				// If the context is canceled, return from the goroutine
				return
			case <-time.Tick(seconds * time.Second):
				// Fetching data to input into image
				cpuTemp, _, cpuFan, _, err := sensors.GetSensorData()
				if err != nil {
					log.Fatal(err)
				}

				// Get CPU Load
				cpuLoad, err := sensors.GetCPULoad()
				if err != nil {
					continue // Skip this cycle if there was an error fetching the data
				}

				// Get network speed
				netSpeed, err := sensors.GetNetworkSpeed()
				if err != nil {
					continue // Skip this cycle if there was an error fetching the data
				}

				// Generate frames for single row text
				frames := make([]*frame.Frame, numFrames)
				for i := 0; i < numFrames; i++ {
					f := frame.NewFrame(frame.Width, frame.Height, fontSize, fontFile)
					f.DrawText("      "+cpuTemp, 1)
					if cpuFan != "0 RPM" {
						f.DrawText("   "+cpuFan, 3)
					} else if netSpeed != "0 B/s" {
						f.DrawText("  "+netSpeed, 3)
					} else {
						// If there is no fan speed or network speed, show the CPU load
						str := fmt.Sprintf("   CPU %d%%", cpuLoad)
						f.DrawText(str, 3)
					}
					f.DrawProgressBar(cpuLoad, 4)
					frames[i] = f
				}

				// Append output to an image
				err = gifcreator.SaveGif(fileName, frames)
				if err != nil {
					continue // Skip this cycle if there was an error saving the GIF
				}

				// Append to anime matrix
				utils.Display(fileName)
			}
		}
	}()

	// Handle quit event from systray
	go func() {
		<-quitCh
		cancel()
	}()

	// Set up the systray menu
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Main loop
	for {
		select {
		case <-mQuit.ClickedCh:
			// Quit event received, send a signal to the goroutine to cancel the context
			quitCh <- struct{}{}
			utils.DisableAnime()
			systray.Quit()
			return
		}
	}
}
