package api

import (
	"MSAgent/collectors"
	"MSAgent/constants"
	"MSAgent/models"
	"MSAgent/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/shirou/gopsutil/v3/host"
	"gorm.io/gorm"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"time"
)

var order = "timestamp desc"

type AgentAPIHandler struct {
	DB *gorm.DB
}

// AgentHealth provides a comprehensive status of the agent and the system.
func (agent *AgentAPIHandler) agentHealth(c *fiber.Ctx) error {

	var cpuMetrics []models.CPUInfo
	var memMetrics []models.MemoryInfo
	var diskMetrics []models.DiskInfo

	agent.DB.Order(order).Limit(10).Find(&cpuMetrics)
	agent.DB.Order(order).Limit(10).Find(&memMetrics)

	//query := `
	//SELECT *
	//FROM disk_infos sd
	//INNER JOIN (
	//	SELECT "AgentID", "Drive", MAX("UpdatedAt") AS MaxUpdated
	//	FROM "SystemDiskData"
	//	GROUP BY "AgentID", "Drive"
	//) latest
	//ON sd."AgentID" = latest."AgentID"
	//AND sd."Drive" = latest."Drive"
	//AND sd."UpdatedAt" = latest.MaxUpdated
	//LIMIT 10;
	//`
	//err := agent.DB.Raw(query).Scan(&diskMetrics).Error
	agent.DB.Order(order).Limit(10).Find(&diskMetrics)

	// Map metrics to desired [timestamp, usage] format
	// Map CPU metrics
	var cpu [][]interface{}
	for _, m := range cpuMetrics {
		cpu = append(cpu, []interface{}{
			m.Timestamp.UnixMilli(),
			math.Round(m.UsagePercent*100) / 100,
		})
	}

	// Map Memory metrics
	var memory [][]interface{}
	for _, m := range memMetrics {
		memory = append(memory, []interface{}{
			m.Timestamp.UnixMilli(),
			math.Round(m.Percentage*100) / 100,
		})
	}

	// Map Disk metrics
	//var disk [][]interface{}
	//for _, m := range diskMetrics {
	//	var usageStr string
	//	if m.TotalSize > 0 {
	//		percent := (float64(m.Used) / float64(m.TotalSize)) * 100
	//		usageStr = fmt.Sprintf("%.2f%%", percent)
	//	} else {
	//		usageStr = "N/A"
	//	}
	//
	//	disk = append(disk, )
	//}

	bootTimestamp, err := host.BootTime()
	var uptimeString string
	if err != nil {
		log.Printf("Could not get boot time: %v", err)
		uptimeString = "N/A"
	} else {
		uptimeDuration := time.Since(time.Unix(int64(bootTimestamp), 0))

		// Format to a more readable string like "72h3m4s"
		uptimeString = time.Now().Add(-uptimeDuration).Format("2006-01-02 15:04:05")
		//uptimeString = uptimeDuration.Truncate(time.Second).String()
	}

	hostname, _ := os.Hostname()

	var lastSync *time.Time
	if len(cpuMetrics) > 0 {
		lastSync = &cpuMetrics[0].Timestamp
	}

	//agentInfo := models.AgentInfo{
	//	Version:    constants.AgentVersion,
	//	AgentID:    constants.AgentID,
	//	IPAddress:  utils.GetHostIP(),
	//	Name:       hostname, // Use hostname for the agent's name
	//	OS:         runtime.GOOS,
	//	LastSync:   lastSync, // Use the timestamp of the last collected metric
	//	SDKVersion: runtime.Version(),
	//}

	// 4. Construct final response
	//response := models.HealthStatus{
	//	SystemInfo: models.SystemInfo{
	//		CPU:    cpuMetrics,
	//		Memory: memMetrics,
	//		Disk:   diskMetrics,
	//	},
	//	Uptime:    uptimeString,
	//	AgentInfo: agentInfo,
	//}

	response := fiber.Map{
		"status": "ok",
		"systemInfo": fiber.Map{
			"cpu":    cpu,
			"memory": memory,
			"disk":   diskMetrics,
		},
		"uptime": uptimeString,
		"agent_info": fiber.Map{
			"version":    constants.AgentVersion,
			"agent_id":   constants.AgentID,
			"IPAddress":  utils.GetHostIP(),
			"name":       hostname,
			"os":         runtime.GOOS,
			"lastSync":   lastSync,
			"SDKVersion": runtime.Version(),
		},
	}

	return c.JSON(response)
}

func (agent *AgentAPIHandler) agentMetricCleanUp(c *fiber.Ctx) error {
	daysStr := c.Query("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	agent.DB.Delete(&models.CPUInfo{}, "days > ?", days)

	return c.JSON(fiber.Map{
		"Success": true,
	})
}

func (agent *AgentAPIHandler) agentTopProcesses(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}

	d, err := collectors.CollectTopProcessesInfo(limit)

	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"Data":    d,
	})
}

// SetupRoutes defines the API endpoints.
func (agent *AgentAPIHandler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/agent", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api.Get("/agent/health", agent.agentHealth)

	api.Get("/agent/sync_complete", agent.agentMetricCleanUp)

	api.Get("/agent/resource-usage", agent.agentTopProcesses)

}
