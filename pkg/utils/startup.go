package utils

import (
	"context"
	"go-anime-matrix-io/pkg/frame"
	"go-anime-matrix-io/pkg/gifcreator"
	"go-anime-matrix-io/pkg/sensors"
	"go-anime-matrix-io/static"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	fileName  = "out.gif" // Name of the output file
	fontSize  = 7
	numFrames = 30 // Number of frames for the animations
	seconds   = 1  // Refresh rate
)

func Startup() context.CancelFunc {
	// Check for necessary stuff are installed or whatnot
	err := checkCommands()
	if err != nil {
		log.Println("Required command is missing:", err)
		return nil
	}

	// Clear & enable matrix display via asusctl
	EnableAnime()

	// Create a ticker that triggers every X seconds
	ticker := time.NewTicker(seconds * time.Second)

	// Make sure anime is disabled when the program ends
	defer DisableAnime()

	// Make sure anime is disabled when the program crashes
	defer HandleCrash()

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
				// If the context is canceled, return from the goroutine
				return
			case <-ticker.C:
				// 1. Fetching data to input into image
				cpuTemp, _, cpuFan, _, err := sensors.GetSensorData()
				if err != nil {
					log.Fatal(err)
				}

				// 2. Get CPU Load
				cpuLoad, err := sensors.GetCPULoad()
				if err != nil {
					continue // Skip this cycle if there was an error fetching the data
				}

				// 3. Get network speed
				netSpeed, err := sensors.GetNetworkSpeed()
				if err != nil {
					continue // Skip this cycle if there was an error fetching the data
				}

				// Generate frames for single row text
				frames := make([]*frame.Frame, numFrames)
				for i := 0; i < numFrames; i++ {
					f := frame.NewFrame(frame.Width, frame.Height, fontSize, static.FontFile)
					f.DrawText("      "+cpuTemp, 1)
					if cpuFan != "0 RPM" {
						f.DrawText("   "+cpuFan, 3)
					} else {
						f.DrawText("     "+netSpeed, 3)
					}
					f.DrawProgressBar(cpuLoad, 4)
					frames[i] = f
				}

				// Append output to an img
				err = gifcreator.SaveGif(fileName, frames)
				if err != nil {
					continue // Skip this cycle if there was an error saving the GIF
				}

				// Append to anime matrix
				Display(fileName)
			}
		}
	}()

	// Wait for an interrupt signal
	<-sigs

	// When the signal is received, cancel the context
	cancel()

	// Wait a bit to allow the goroutines to clean up
	time.Sleep(time.Second)

	return cancel
}

// HandleCrash is run to disable anime matrix after a recovery
func HandleCrash() {
	if r := recover(); r != nil {
		log.Println("Recovered from crash, disabling anime:", r)
		DisableAnime()
	}
}
