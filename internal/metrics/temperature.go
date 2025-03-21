package metrics

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/blccming/goSMA/internal/helpers"
)

type TemperatureMetrics struct {
	CPU    float32            `json:"cpu"`
	GPU    float32            `json:"gpu"`
	Drives map[string]float32 `json:"drives"`
}

// Fetches every hwmon device from /sys/class/hwmon
//
// Returns:
//   - map in the format of hwmon0..*:name (name being /sys/class/hwmon/hwmon*/name)
func getHwonDevices() map[string]string {
	devices := make(map[string]string)
	hwmonList := helpers.ListAsString("/sys/class/hwmon")

	for _, hwmon := range hwmonList {
		deviceType := string(helpers.ReadFile("/sys/class/hwmon/" + hwmon + "/name"))
		devices[hwmon] = strings.Trim(deviceType, "\n")
	}

	return devices
}

// Parses temperatures from all relevant hwmon devices
//
// Returns:
//   - map in the format of device:temperature
func getTemperatures() map[string]float32 {
	temps := make(map[string]float32)
	devices := getHwonDevices()

	cpuNames := []string{"cpu_thermal", "coretemp", "k10temp"}
	nvmeCounter := 0
	deviceName := ""

	for hwmon, device := range devices {
		if !slices.Contains([]string{"nvme", "amdgpu"}, device) && !slices.Contains(cpuNames, device) {
			continue
		}

		// parse temperature value
		tempString := string(helpers.ReadFile("/sys/class/hwmon/" + hwmon + "/temp1_input"))
		tempStringTimmed := strings.Trim(tempString, "\n")
		tempInt, err := strconv.Atoi(tempStringTimmed)
		if err != nil {
			helpers.LogError(fmt.Errorf("getTemperatures(): Error while converting string to float: %w", err))
		}
		tempFloat := float32(tempInt) / 1000

		// choose device type
		if device == "nvme" {
			deviceName = "nvme" + strconv.Itoa(nvmeCounter)
			nvmeCounter++
		} else if slices.Contains(cpuNames, device) {
			deviceName = "CPU"
		} else if device == "amdgpu" {
			deviceName = "GPU"
		}

		temps[deviceName] = tempFloat
	}

	return temps
}

// Fetch temperature info from host
//
// Returns:
//   - Struct including temperature from cpu, gpu and nvme drives
func Temperature() TemperatureMetrics {
	tempsData := getTemperatures()

	temperature := TemperatureMetrics{
		CPU:    tempsData["CPU"],
		GPU:    tempsData["GPU"],
		Drives: make(map[string]float32),
	}

	for device, temp := range tempsData {
		if strings.HasPrefix(device, "nvme") {
			temperature.Drives[device] = temp
		}
	}

	return temperature
}
