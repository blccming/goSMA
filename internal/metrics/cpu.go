package metrics

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blccming/goSMA/internal/helpers"
)

type CpuMetrics struct {
	Model   string  `json:"model"`
	Cores   int     `json:"cores"`
	Threads int     `json:"threads"`
	Usage   float32 `json:"usage"`
}

// Reads the current CPU time from /proc/stat file
//
// Returns:
//   - Total CPU time
//   - Idle CPU time
func getCpuTime() (int, int) {
	procStat := helpers.ReadFileAsString("/proc/stat")
	statFirstLine := strings.Split(procStat, "\n")[0]

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

// Reads general information about the cpu from /proc/cpuinfo
//
// Returns:
//   - Model (e.g. "AMD Ryzen 7 5800X 8-Core Processor")
//   - Core count
//   - Thread count
func getCpuInfo() (string, int, int) {
	procCpuinfo := helpers.ReadFileAsString("/proc/cpuinfo")
	procCpuInfoParts := strings.Split(procCpuinfo, "\n")

	cpuinfoLineModel := procCpuInfoParts[4]
	cpuinfoLineThreads := procCpuInfoParts[10]
	cpuinfoLineCores := procCpuInfoParts[12]

	model := strings.Split(cpuinfoLineModel, ": ")[1]
	threadsString := strings.Split(cpuinfoLineThreads, ": ")[1]
	coresString := strings.Split(cpuinfoLineCores, ": ")[1]

	threadsInt, err := strconv.Atoi(threadsString)
	if err != nil {
		helpers.LogError(fmt.Errorf("getCpuInfo(): Error while converting string to int: %w", err))
		return model, 0, 0
	}

	coresInt, err := strconv.Atoi(coresString)
	if err != nil {
		helpers.LogError(fmt.Errorf("getCpuInfo(): Error while converting string to int: %w", err))
		return model, threadsInt, 0
	}

	return model, threadsInt, coresInt
}

// Fetch CPU info from host
//
// Returns:
//   - Struct including the CPU model, core count, thread count and usage percentage
func CPU() CpuMetrics {
	model, threads, cores := getCpuInfo()

	cpu := CpuMetrics{
		Model:   model,
		Cores:   cores,
		Threads: threads,
		Usage:   getCpuUsage(),
	}

	return cpu
}
