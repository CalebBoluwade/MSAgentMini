package collectors

import (
	"MSAgent/models"
	psnet "github.com/shirou/gopsutil/v3/net"
	"gorm.io/gorm"
	"log"
	"time"
)

func CollectNetworkInfo(db *gorm.DB) {
	ioCounters, err := psnet.IOCounters(true)
	if err != nil {
		log.Printf("Error getting network usage: %v", err)
		return
	}

	for _, counter := range ioCounters {
		usage := models.NetworkInfo{
			Metric:      models.Metric{Timestamp: time.Now().UTC()},
			Name:        counter.Name,
			BytesSent:   counter.BytesSent,
			BytesRecv:   counter.BytesRecv,
			PacketsSent: counter.PacketsSent,
			PacketsRecv: counter.PacketsRecv,
		}

		log.Println("Network", usage)
		if err := db.Create(&usage).Error; err != nil {
			log.Printf("Error saving network usage for interface %s: %v", counter.Name, err)
		}
	}
}
