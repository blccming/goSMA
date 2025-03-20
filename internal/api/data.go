package api

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/blccming/goSMA/internal/helpers"
	"github.com/blccming/goSMA/internal/metrics"
)

type EndpointData struct {
	CPU     metrics.CpuMetrics     `json:"cpu"`
	Memory  metrics.MemoryMetrics  `json:"memory"`
	Network metrics.NetworkMetrics `json:"network"`
	System  metrics.SystemMetrics  `json:"system"`
}

var data EndpointData

// Updates the EndpiontData object (fetches every metric)
func updateData() {
	data = EndpointData{
		CPU:     metrics.CPU(),
		Memory:  metrics.Memory(),
		Network: metrics.Network(),
		System:  metrics.System(),
	}
}

// Runs updateData() every three seconds after initial execution of updateData()
func StartUpdating() {
	updateData() // Initial exection so actual data can be returned asap

	// Get update intervall from environment
	updateIntervallString := os.Getenv("UPDATE_INTERVALL")
	if updateIntervallString == "" {
		updateIntervallString = "1.0"
	}

	updateIntervall, err := strconv.ParseFloat(updateIntervallString, 32)
	if err != nil {
		helpers.LogError(fmt.Errorf("StartUpdating(): Error while converting string to float: %w. Using 1.0 seconds as update intervall.", err))
		updateIntervall = 1.0
	}

	// Starts ticker
	ticker := time.NewTicker(time.Duration(updateIntervall*1000) * time.Millisecond)

	// Stops ticker when function exits
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go updateData() // Utilize go coroutine to make update interval more variable (some metric fetching function utilize sleep times)
		}
	}
}

// Returns:
//   - EndpointData object
func getEndpointData() EndpointData {
	return data
}
