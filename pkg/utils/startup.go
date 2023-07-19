package utils

import (
	"context"
	"go-anime-matrix-io/internal/models"
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
	numFrames = 1 // Number of frames for the animations
	seconds   = 1 // Refresh rate
)

func Startup(settings *models.AppSettings) context.CancelFunc {
	// Check for necessary stuff are installed or not
	err := checkCommands()
	if err != nil {
		// TODO: Add a popup to the GUI
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
				// If the ticker triggers, generate a new GIF

				var frames []*frame.Frame

				if settings.Mode == "System mode" {
					var cpuTemp, _, fanSpeed, gpuFan string
					var cpuLoadOrRAMUsed int

					// Fetching CPU temp
					if settings.ShowCPUTemp {
						cpuTemp, _, fanSpeed, gpuFan, err = sensors.GetSensorData()
						if err != nil {
							log.Fatal(err)
						}
					} else {
						cpuTemp = ""
					}

					// Select which sensor to use
					if settings.CPUFanSpeed == "GPU Fan Speed" {
						fanSpeed = gpuFan
					}

					switch settings.CPUFanSpeed {
					case "CPU Fan Speed":
						// Defaults to CPU Fan Speed
					case "GPU Fan Speed":
						fanSpeed = gpuFan
					case "Average Fan Speeds":
						fanSpeed = sensors.GetAverageFanSpeed(fanSpeed, gpuFan)
					}

					// Fetching CPU Load or RAM Use
					if settings.CPULoadOrRAMUse == "CPU Load" {
						// Get CPU Load
						cpuLoadOrRAMUsed, err = sensors.GetCPULoad()
						if err != nil {
							continue // Skip this cycle if there was an error fetching the data
						}
					} else {
						// Get RAM Use
						cpuLoadOrRAMUsed, err = sensors.GetRAMUsage()
						if err != nil {
							continue // Skip this cycle if there was an error fetching the data
						}
					}

					// Get network speed
					netSpeed, err := sensors.GetNetworkSpeed()
					if err != nil {
						continue // Skip this cycle if there was an error fetching the data
					}

					// Generate frames for single row text
					frames = make([]*frame.Frame, numFrames)
					for i := 0; i < numFrames; i++ {
						f := frame.NewFrame(frame.Width, frame.Height, fontSize, static.FontFile)
						f.DrawText("      "+cpuTemp, 1)
						if fanSpeed != "" {
							f.DrawText("   "+fanSpeed, 3)
						} else {
							f.DrawText("     "+netSpeed, 3)
						}
						f.DrawProgressBar(cpuLoadOrRAMUsed, 4)
						frames[i] = f
					}
				}

				// TODO: Add audio mode
				if settings.Mode == "Audio mode" {
					// If demo mode is enabled, generate random data
					if settings.EqualizerDemo {
						frames = sensors.GenerateEqualizerFrames(10, 12)
					} else {
						// TODO: Get audio data in frames
						// TODO: Get the currently playing song's name
					}
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
