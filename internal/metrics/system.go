package metrics

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/blccming/goSMA/internal/helpers"
)

type SystemMetrics struct {
	Hostname     string  `json:"hostname"`
	Distribution string  `json:"distribution"`
	LinuxVersion string  `json:"linux_version"`
	Uptime       float64 `json:"uptime"`
	State        string  `json:"state"`
}

// Get hostname from either environment variable or /etc/hostname
//
// Returns:
//   - Hostname
func getHostname() string {
	envHostname := os.Getenv("HOSTNAME")
	if envHostname != "" {
		return envHostname
	}

	etcHostname := helpers.ReadFileAsString("/etc/hostname")

	return etcHostname
}

// Fetch PRETTY_NAME from /etc/os-release
//
// Returns:
//   - Pretty name (e.g. "Fedora Linux 41 (Workstation Edition)")
func getDistribution() string {
	etcRelease := helpers.ReadFileAsString("/etc/os-release")

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
	procVersion := helpers.ReadFileAsString("/proc/version")
	linuxVersion := strings.Split(string(procVersion), " ")[2]

	return linuxVersion
}

// Fetch uptime from /proc/uptime
//
// Returns:
//   - Uptime in seconds
func getUptime() float64 {
	procUptime := helpers.ReadFileAsString("/proc/uptime")
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
//   - Struct including distribution, linux version and state (online)
func System() SystemMetrics {
	system := SystemMetrics{
		Hostname:     getHostname(),
		Distribution: getDistribution(),
		LinuxVersion: getLinuxVersion(),
		Uptime:       getUptime(),
		State:        "online",
	}

	return system
}
