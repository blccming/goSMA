package metrics

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blccming/goSMA/internal/helpers"
)

type cpuInfo struct {
	Model   string  `json:"model"`
	Cores   int8    `json:"cores"`
	Threads int8    `json:"threads"`
	Usage   float32 `json:"usage"`
}

// Reads the current CPU time from /proc/stat file
//
// Returns:
//   - Total CPU time
//   - Idle CPU time
func getCpuTime() (int, int) {
	stat := string(helpers.ReadFile("/proc/stat")[:])
	statFirstLine := strings.Split(stat, "\n")[0]

	cpuTimesString := strings.Split(statFirstLine, " ")[2:] // remove "cpu" and extra space after "cpu"

	cpuTimesInt := make([]int, len(cpuTimesString)) // create cpuTimesInt with length of cpuTimesString

	var err error
	for i, cpuTimeString := range cpuTimesString {
		cpuTimesInt[i], err = strconv.Atoi(cpuTimeString)
		if err != nil {
			helpers.LogError(fmt.Errorf("getCpuTime(): Error while converting string to int: %w", err))
			return 0, 0
		}
	}

	var idleTime, totalTime int
	for i, cpuTime := range cpuTimesInt {
		if i == 3 {
			idleTime = cpuTime // idle time is the fourth element of cpuTimesInt[]
		}
		totalTime += cpuTime
	}

	return totalTime, idleTime
}

// Calculates current CPU usage
//
// Returns:
//   - CPU usage percentage
func getCpuUsage() float32 {
	totalT1, idleT1 := getCpuTime()
	time.Sleep(1 * time.Second)
	totalT2, idleT2 := getCpuTime()

	usage := float32(((totalT2 - totalT1) - (idleT2 - idleT1))) / float32(totalT2-totalT1) * 100
	return usage
}

// Fetch CPU info from host
//
// Returns:
//   - JSON string including the CPU model, core count, thread count and usage percentage
func CPU() string {
	cpu := cpuInfo{
		Model:   "",
		Cores:   0,
		Threads: 0,
		Usage:   getCpuUsage(),
	}

	jsonData, err := json.Marshal(cpu)
	if err != nil {
		helpers.LogError(fmt.Errorf("CPU(): Error marshalign to JSON: %w", err))
		return ""
	}

	return string(jsonData)
}
