package sensors

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/elastic/gosigar"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/net"
	"go-anime-matrix-io/pkg/frame"
	"go-anime-matrix-io/static"
	"math"
	"math/rand"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GetSensorData returns and structures the output from the lm-sensors package
func GetSensorData() (cpuTemp string, gpuTemp string, cpuFan string, gpuFan string, err error) {
	cmd := exec.Command("sensors")
	out, err := cmd.Output()
	if err != nil {
		return "", "", "", "", err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	cpuTempRegexp := regexp.MustCompile(`Tctl:\s+\+([\d.]+)°C`)
	gpuTempRegexp := regexp.MustCompile(`edge:\s+\+([\d.]+)°C`)
	cpuFanRegexp := regexp.MustCompile(`cpu_fan:\s+(\d+) RPM`)
	gpuFanRegexp := regexp.MustCompile(`gpu_fan:\s+(\d+) RPM`)

	for scanner.Scan() {
		line := scanner.Text()
		if cpuTempRegexp.MatchString(line) {
			cpuTemp = cpuTempRegexp.FindStringSubmatch(line)[1]
			var flt float64
			flt, err = strconv.ParseFloat(cpuTemp, 64)
			if err != nil {
				return
			}

			// Round the float up to the nearest whole number
			rounded := math.Ceil(flt)

			// Convert the float to an integer
			cpuTemp = fmt.Sprintf("%v°C", int(rounded))
		}

		if gpuTempRegexp.MatchString(line) {
			gpuTemp = gpuTempRegexp.FindStringSubmatch(line)[1]
			var flt float64
			flt, err = strconv.ParseFloat(gpuTemp, 64)
			if err != nil {
				return
			}

			// Round the float up to the nearest whole number
			rounded := math.Ceil(flt)

			// Convert the float to an integer
			gpuTemp = fmt.Sprintf("%v°C", int(rounded))
		}

		if cpuFanRegexp.MatchString(line) {
			cpuFan = cpuFanRegexp.FindStringSubmatch(line)[1]
			cpuFan = fmt.Sprintf("%s RPM", cpuFan)
		}

		if gpuFanRegexp.MatchString(line) {
			gpuFan = gpuFanRegexp.FindStringSubmatch(line)[1]
			gpuFan = fmt.Sprintf("%s RPM", gpuFan)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", "", "", err
	}

	return cpuTemp, gpuTemp, cpuFan, gpuFan, nil
}

// GetCPULoad returns the current CPU load
func GetCPULoad() (load int, err error) {
	percent, err := cpu.Percent(time.Second, false)

	// Round the float up to the nearest whole number
	rounded := math.Ceil(percent[0])
	return int(rounded), err
}

// GetRAMUsage returns the current RAM usage in percentage
func GetRAMUsage() (usage int, err error) {
	var mem gosigar.Mem
	if err = mem.Get(); err != nil {
		return 0, err
	}

	// RAM usage is calculated as ActualUsed / Total * 100
	usage = int(float64(mem.ActualUsed) / float64(mem.Total) * 100)
	return usage, nil
}

// GenerateEqualizerFrames returns a string of equalizer bars based on the input
func GenerateEqualizerFrames(length, maxHeight int) []*frame.Frame {
	frames := make([]*frame.Frame, 10)

	// Seed the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 10; i++ {
		f := frame.NewFrame(frame.Width, frame.Height, 6, static.FontFile)

		// Generate equalizer blocks for each frame
		for h := 1; h <= maxHeight; h++ {
			equalizerLine := ""
			for j := 0; j < length; j++ {
				// Generate a random height for the current block
				blockHeight := rand.Intn(maxHeight) + 1

				if j%(maxHeight+1) <= maxHeight-blockHeight {
					equalizerLine += "█" // Full block character
				} else {
					equalizerLine += " " // Space character
				}
			}
			f.DrawText(equalizerLine, h)
		}

		frames[i] = f
	}

	return frames
}

// GetNetworkSpeed returns the current network speed
func GetNetworkSpeed() (string, error) {
	// Get initial network stats
	statsStart, err := net.IOCounters(true)
	if err != nil {
		return "", err
	}

	// Sleep for a second
	time.Sleep(1 * time.Second)

	// Get network stats again
	statsEnd, err := net.IOCounters(true)
	if err != nil {
		return "", err
	}

	// Calculate the difference in bytes sent and received
	bytesRecv := statsEnd[0].BytesRecv - statsStart[0].BytesRecv

	// Format bytes to appropriate units
	recv := formatBytes(bytesRecv)

	// Return as a string
	return fmt.Sprintf("%s", recv), nil
}

// Helpers

// GetAverageFanSpeed returns the average temperature of the CPU and GPU
func GetAverageFanSpeed(cpuFan, gpuFan string) string {
	cpuFanSpeed, err := strconv.Atoi(strings.TrimSuffix(cpuFan, " RPM"))
	if err != nil {
		return "0 RPM"
	}

	gpuFanSpeed, err := strconv.Atoi(strings.TrimSuffix(gpuFan, " RPM"))
	if err != nil {
		return "0 RPM"
	}

	average := (cpuFanSpeed + gpuFanSpeed) / 2
	return fmt.Sprintf("%d RPM", average)
}

// formatBytes formats bytes to a human-readable format
func formatBytes(bytes uint64) string {
	const (
		_         = iota
		KB uint64 = 1 << (10 * iota)
		MB
		GB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB/s", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB/s", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB/s", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B/s", bytes)
	}
}
