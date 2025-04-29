package metrics

import (
	"golang.org/x/sys/unix"
)

type StorageMetrics struct {
	Device  string  `json:"device"`
	Size    uint64  `json:"size"`
	Used    uint64  `json:"used"`
	UsedPct float32 `json:"usage"`
}

// Uses the unix stdlib to parse root filesystem stat (from "/")
//
// Returns:
//   - Total size in kB
//   - Used size in kB
//   - Usage in percent
func parseRootStat() (total uint64, used uint64, usedpct float32) {
	var stat unix.Statfs_t
	unix.Statfs("/", &stat)

	total = stat.Blocks * uint64(stat.Bsize) // Total size in bytes
	free := stat.Bfree * uint64(stat.Bsize)
	used = total - free

	usedpct = float32(used) / float32(total) * 100

	return total, used, usedpct
}

// Fetch storage info from host
//
// Returns:
//   - Struct including the storage device, storage size, used storage (absolute) and used storage (relative)
func Storage() StorageMetrics {
	total, used, usedpct := parseRootStat()

	storage := StorageMetrics{
		Device:  "rootfs",
		Size:    total,
		Used:    used,
		UsedPct: usedpct,
	}

	return storage
}

// will be used at a later point in time
//
// type StorageDevice struct {
// 	Mountpoint string  `json:"mountpoint"`
// 	SizeKb     int     `json:"size"`
// 	UsageKb    int     `json:"usage_absolute"`
// 	UsagePct   float32 `json:"usage_relative"`
// 	IoRead     int     `json:"io_read"`
// 	IoWrite    int     `json:"io_write"`
// }

// type StorageMetrics struct {
// 	Devices []StorageDevice `json:"devices"`
// }

// func parseProcDiskstats() map[string][]int {
// 	diskstats := make(map[string][]int)

// 	procDiskstats := helpers.ReadFileAsString("/proc/diskstats")
// 	procDiskstatsTrimmed := strings.TrimSpace(procDiskstats)
// 	procDiskstatsSplitted := strings.Split(procDiskstatsTrimmed, "\n")

// 	for _, line := range procDiskstatsSplitted {
// 		sections := strings.Split(line, " ")
// 		device := ""
// 		for _, s := range sections {
// 			if s == " " {
// 				continue
// 			}

// 			// get first rune of section to check if its a letter
// 			r, _ := utf8.DecodeRuneInString(s)
// 			if unicode.IsLetter(r) {
// 				device = s // if it has a letter, set section as device name
// 			}

// 			// skip first sections until device name has been set (third section or so in /proc/diskstats)
// 			if device == "" {
// 				continue
// 			}

// 			// if section is not the device name, convert it as int and append the value to map
// 			if s != device {
// 				valueInt, err := strconv.Atoi(s)
// 				if err != nil {
// 					helpers.LogError(fmt.Errorf("parseProcDiskstats(): Error while converting string to int: %w", err))
// 				}

// 				diskstats[device] = append(diskstats[device], valueInt)
// 			}
// 		}
// 	}

// 	return diskstats
// }
