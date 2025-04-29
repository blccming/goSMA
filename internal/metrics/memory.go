package metrics

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blccming/goSMA/internal/helpers"
)

type MemoryMetrics struct {
	Total      int `json:"total"`
	Used       int `json:"used"`
	Cached     int `json:"cached"`
	Free       int `json:"free"`
	Available  int `json:"available"`
	SwapTotal  int `json:"swap_total"`
	SwapCached int `json:"swap_cached"`
	SwapFree   int `json:"swap_free"`
}

// Parses /proc/meminfo and filters out only the relevant lines
//
// Returns
//   - A map: memory metric name:amount in kilobytes
func parseProcMeminfo() map[string]int {
	procMeminfo := helpers.ReadFileAsString("/proc/meminfo")

	splittedMeminfo := strings.Split(procMeminfo, "\n")

	filteredMeminfo := make(map[string]int)
	filters := []string{"MemTotal", "MemFree", "MemAvailable", "Cached", "SwapTotal", "SwapFree", "SwapCached"}

	for _, line := range splittedMeminfo {
		for _, filter := range filters {
			if strings.HasPrefix(line, filter) {
				data := strings.Split(line, ":")[1]
				dataNoSuffix := strings.TrimSuffix(data, "kB")
				dataNoSuffixTrimmed := strings.Trim(dataNoSuffix, " ")

				dataInt, err := strconv.Atoi(dataNoSuffixTrimmed)
				if err != nil {
					helpers.LogError(fmt.Errorf("parseProcMeminfo(): Error while converting string to int: %w", err))
				}

				filteredMeminfo[filter] = dataInt
			}
		}
	}

	return filteredMeminfo
}

// Fetch memory info from host
//
// Returns:
//   - Struct including total, used, cached, free and available memory as well as the total, cached and free amount of swap
func Memory() MemoryMetrics {
	meminfo := parseProcMeminfo()

	memory := MemoryMetrics{
		Total:      meminfo["MemTotal"],
		Used:       meminfo["MemTotal"] - meminfo["MemAvailable"],
		Cached:     meminfo["Cached"],
		Free:       meminfo["MemFree"],
		Available:  meminfo["MemAvailable"],
		SwapTotal:  meminfo["SwapTotal"],
		SwapCached: meminfo["SwapCached"],
		SwapFree:   meminfo["SwapFree"],
	}

	return memory
}
