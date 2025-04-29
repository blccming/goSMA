package metrics

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blccming/goSMA/internal/helpers"
)

type NetworkMetrics struct {
	UsedInterfaces []string `json:"used_interfaces"`
	UsedBandwithRX uint64   `json:"rx"`
	UsedBandwithTX uint64   `json:"tx"`
}

// Parses /proc/net/dev and filters out any irrelevant lines, only keeps the lines with interfaces that have desired prefixes
//
// Returns:
//   - List with every relevant line from /proc/net/dev
func parseNetDev() []string {
	interfacePrefixes := []string{"enp", "wlp", "eth", "wlan"}

	procNetDev := helpers.ReadFileAsString("/proc/net/dev")
	procNetDevLines := strings.Split(procNetDev, "\n")

	// Filter out any line with interfaces that are not in interfacePrefixes
	var filteredLines []string
	for _, line := range procNetDevLines {
		for _, prefix := range interfacePrefixes {
			if strings.HasPrefix(line, prefix) {
				filteredLines = append(filteredLines, line)
				break
			}
		}
	}

	// Join consecutive whitespaces to one
	for i, line := range filteredLines {
		filteredLines[i] = strings.Join(strings.Fields(line), " ")
	}

	return filteredLines
}

// Filters out used interfaces from NetDev
//
// Returns:
//   - List with all used interfaces
func getInterfaces() []string {
	NetDev := parseNetDev()

	var interfaces []string
	for _, line := range NetDev {
		splittedLine := strings.Split(line, " ")

		// Reject interface if it has received 0 bytes so far
		if splittedLine[1] != "0" {
			interfacePart := strings.ReplaceAll(splittedLine[0], ":", "")
			interfaces = append(interfaces, interfacePart)
		}
	}

	return interfaces
}

// Aggregates total received and transmitted bytes from multiple interfaces (valid interfaces defined in parseNetDev())
//
// Returns:
//   - Total received bytes (limit: 1.844674407×10¹⁹ bytes)
//   - Total transmittes bytes (limit: 1.844674407×10¹⁹ bytes)
func getTotalBytes() (uint64, uint64) {
	NetDev := parseNetDev()

	// Go through every interface in NetDev and add total received/transmitted byytes to rxBytes/txBytes
	var rxBytes, txBytes uint64
	for _, line := range NetDev {
		lineSplitted := strings.Split(line, " ")

		// Filter out information for total received and total transmitted bytes, see /proc/net/dev for indexes
		rxBytesString := lineSplitted[1]
		txBytesString := lineSplitted[9]

		rxBytesInt, err := strconv.ParseUint(rxBytesString, 10, 64)
		if err != nil {
			helpers.LogError(fmt.Errorf("getTotalBytes(): Error while converting string to int: %w", err))
		}
		rxBytes += rxBytesInt

		txBytesInt, err := strconv.ParseUint(txBytesString, 10, 64)
		if err != nil {
			helpers.LogError(fmt.Errorf("getTotalBytes(): Error while converting string to int: %w", err))
		}
		txBytes += txBytesInt
	}

	return rxBytes, txBytes
}

// Calculates currently used bandwith (rx, tx)
//
// Returns:
//   - Received bytes per second (limit: 1.844674407×10¹⁹ bytes per second)
//   - Transmitted bytes per second (limit: 1.844674407×10¹⁹ bytes per second)
func getNetworkUsage() (uint64, uint64) {
	totalRx1, totalTx1 := getTotalBytes()
	time.Sleep(1 * time.Second)
	totalRx2, totalTx2 := getTotalBytes()

	rx := totalRx2 - totalRx1
	tx := totalTx2 - totalTx1

	return rx, tx
}

// Fetch network info from host
//
// Returns:
//   - Struct including main network interfaces, receiving bytes per second, transmitting bytes per second
func Network() NetworkMetrics {
	rx, tx := getNetworkUsage()

	network := NetworkMetrics{
		UsedInterfaces: getInterfaces(),
		UsedBandwithRX: rx,
		UsedBandwithTX: tx,
	}

	return network
}
