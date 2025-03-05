package api

import (
	"time"

	"github.com/blccming/goSMA/internal/metrics"
)

type EndpointData struct {
	CPU     metrics.CpuMetrics     `json:"cpu"`
	Network metrics.NetworkMetrics `json:"network"`
	System  metrics.SystemMetrics  `json:"system"`
}

var data EndpointData

// Updates the EndpiontData object (fetches every metric)
func updateData() {
	data = EndpointData{
		CPU:     metrics.CPU(),
		Network: metrics.Network(),
		System:  metrics.System(),
	}
}

// Runs updateData() every three seconds after initial execution of updateData()
func StartUpdating() {
	updateData() // Initial exection so actual data can be returned asap

	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop() // Stops timer when function exits

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
