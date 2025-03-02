package api

import (
	"time"

	"github.com/blccming/goSMA/internal/metrics"
)

type EndpointData struct {
	CPU    metrics.CpuMetrics    `json:"cpu"`
	System metrics.SystemMetrics `json:"system"`
}

var data EndpointData

// Updates the EndpiontData object (fetches every metric)
func updateData() {
	data = EndpointData{
		CPU:    metrics.CPU(),
		System: metrics.System(),
	}
}

// Runs updateData() every three seconds after initial execution of updateData()
func StartUpdating() {
	updateData() // Initial exection so actual data can be returned asap

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop() // Stops timer when function exits

	for {
		select {
		case <-ticker.C:
			updateData()
		}
	}
}

// Returns:
//   - EndpointData object
func getEndpointData() EndpointData {
	return data
}
