package collectors

import (
	"MSAgent/models"
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"strings"
)

func CollectDiskSpace(db *gorm.DB, returnListObj bool) []models.DiskInfo {
	var formattedDiskInfo []models.DiskInfo
	seenDrives := make(map[string]struct{}) // Tracks unique Drive identifiers

	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Disk Partitions Error: %v", err)
		return formattedDiskInfo
	}

	for _, part := range partitions {
		// Filter out unwanted mountpoints/devices
		if strings.Contains(part.Mountpoint, "app") || strings.Contains(part.Mountpoint, "System") || strings.Contains(part.Device, "devfs") && strings.Contains(part.Device, "dev") {
			continue
		}

		slog.Info("dd", "mt", part.Mountpoint, "dv", part.Device, "fs", part.Fstype)

		// Normalize drive name (remove backslashes, trim)
		drive := strings.TrimSpace(strings.ReplaceAll(part.Device, "\\", ""))

		// Deduplicate by device path
		if _, exists := seenDrives[drive]; exists {
			continue
		}
		seenDrives[drive] = struct{}{}

		usage, err := disk.Usage(part.Mountpoint)
		if err != nil || usage.Total == 0 {
			continue
		}

		diskUsage := models.DiskInfo{
			Drive:      drive,
			TotalSize:  usage.Total,
			Free:       usage.Free,
			Used:       usage.Used,
			FormatSize: fmt.Sprintf("%.2f GB", float64(usage.Total)/(1024*1024*1024)),
			FormatFree: fmt.Sprintf("%.2f GB", float64(usage.Free)/(1024*1024*1024)),
		}

		log.Println("DiskInfo", drive, usage.Total, diskUsage)

		if err := db.Create(&diskUsage).Error; err != nil {
			log.Printf("Error saving Disk usage: %v", err)
		}

		formattedDiskInfo = append(formattedDiskInfo, diskUsage)
	}

	return formattedDiskInfo
}
