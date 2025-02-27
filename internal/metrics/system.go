package metrics

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/blccming/goSMA/internal/helpers"
)

type systemMetrics struct {
	Distribution string  `json:"distribution"`
	LinuxVersion string  `json:"linux_version"`
	Uptime       float64 `json:"uptime"`
	State        string  `json:"state"`
}

// Fetch PRETTY_NAME from /etc/os-release
//
// Returns:
//   - Pretty name (e.g. "Fedora Linux 41 (Workstation Edition)")
func getDistribution() string {
	etcRelease := string(helpers.ReadFile("/etc/os-release"))

	etcReleaseParts := strings.Split(etcRelease, "\n")

	for _, part := range etcReleaseParts {
		if strings.HasPrefix(part, "PRETTY_NAME=") {
			prettyName := strings.Split(part, "=")[1]
			prettyName = strings.Trim(prettyName, "\"")
			return prettyName
		}
	}

	return "PRETTY_NAME NOT DEFINED"
}

// Fetch linux version from /proc/version
//
// Returns:
//   - Linux version (e.g. "6.13.4-200.fc41.x86_64")
func getLinuxVersion() string {
	procVersion := string(helpers.ReadFile("/proc/version"))
	linuxVersion := strings.Split(string(procVersion), " ")[2]

	return linuxVersion
}

// Fetch uptime from /proc/uptime
//
// Returns:
//   - Uptime in seconds
func getUptime() float64 {
	procUptime := string(helpers.ReadFile("/proc/uptime"))
	uptimeString := strings.Split(procUptime, " ")[0]

	uptimeFloat, err := strconv.ParseFloat(uptimeString, 64)
	if err != nil {
		helpers.LogError(fmt.Errorf("getUptime(): Error while converting string to float32 %w", err))
	}

	return uptimeFloat
}

// Fetch system info from host
//
// Returns:
//   - JSON string including distribution, linux version and state (online)
func System() string {
	system := systemMetrics{
		Distribution: getDistribution(),
		LinuxVersion: getLinuxVersion(),
		Uptime:       getUptime(),
		State:        "online",
	}

	jsonData, err := json.Marshal(system)
	if err != nil {
		helpers.LogError(fmt.Errorf("System(): Error marshaling to JSON: %w", err))
	}

	return string(jsonData)
}
