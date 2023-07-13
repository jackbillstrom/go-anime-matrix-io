package sensors

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"strconv"
)

// GetSensorData returns and sturctures the output from the lm-sensors package
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


func convertStringToFloat(str string) (float64, error) {
	flt, err := strconv.ParseFloat(str, 64)
	if err != nil {
		// do something with err
	}
	return flt, err // <-- err is in scope here, no error
}