package collectors

import (
	"MSAgent/models"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"gorm.io/gorm"
	"log"
	"os"
	"sort"
	"time"
)

// CollectCPUInfo gets CPU usage and stores it.
func CollectCPUInfo(db *gorm.DB) {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		return
	}

	if len(percentages) > 0 {
		usage := models.CPUInfo{
			Metric:       models.Metric{Timestamp: time.Now().UTC()},
			UsagePercent: percentages[0],
		}

		log.Println("CPUInfo", usage)

		if err := db.Create(&usage).Error; err != nil {
			log.Printf("Error saving CPU usage: %v", err)
		}
	}
}

// CollectMemoryInfo gets memory usage and stores it.
func CollectMemoryInfo(db *gorm.DB) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting memory usage: %v", err)
		return
	}

	// fmt.Sprintf("%.2f GB",
	usage := models.MemoryInfo{
		Metric:     models.Metric{Timestamp: time.Now().UTC()},
		Total:      float64(vm.Total) / (1024 * 1024 * 1024),
		Available:  float64(vm.Available) / (1024 * 1024 * 1024),
		Used:       float64(vm.Used) / (1024 * 1024 * 1024),
		Percentage: vm.UsedPercent,
	}

	log.Println("MemoryInfo", usage)
	if err := db.Create(&usage).Error; err != nil {
		log.Printf("Error saving memory usage: %v", err)
	}
}

func CollectTopProcessesInfo(limit int) ([]models.ProcInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var processList []models.ProcInfo

	for _, p := range procs {
		name, _ := p.Name()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		username, _ := p.Username()
		statuses, _ := p.Status()
		status := ""
		if len(statuses) > 0 {
			status = statuses[0]
		}
		createTimeMillis, _ := p.CreateTime()
		createTime := time.Unix(0, createTimeMillis*int64(time.Millisecond))

		processList = append(processList, models.ProcInfo{
			PID:        p.Pid,
			Name:       name,
			CPUPercent: cpuPercent,
			MemPercent: memPercent,
			User:       username,
			Status:     status,
			CreateTime: createTime,
		})
	}

	// Sort by CPU and memory usage
	sort.Slice(processList, func(i, j int) bool {
		if processList[i].CPUPercent == processList[j].CPUPercent {
			return processList[i].MemPercent > processList[j].MemPercent
		}
		return processList[i].CPUPercent > processList[j].CPUPercent
	})

	if len(processList) > limit {
		processList = processList[:limit]
	}

	printProcessTable(processList)

	return processList, nil
}

func printProcessTable(procs []models.ProcInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"PID", "Name", "Status", "Created", "User", "CPU %", "Memory %"})

	for _, p := range procs {
	_:
		table.Append([]string{
			fmt.Sprintf("%d", p.PID),
			p.Name,
			p.Status,
			p.CreateTime.Format("2006-01-02 15:04:05"),
			p.User,
			fmt.Sprintf("%.2f%%", p.CPUPercent),
			fmt.Sprintf("%.2f%%", p.MemPercent),
		})
	}
_:
	table.Render()
}
