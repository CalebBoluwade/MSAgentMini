package models

import (
	"time"
)

type Metric struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// CPUInfo holds CPU utilization data.
type CPUInfo struct {
	Metric
	UsagePercent float64 `json:"usage_percent"`
}

// MemoryInfo holds memory statistics.
type MemoryInfo struct {
	Metric
	Total      float64 `json:"Total"`
	Available  float64 `json:"Available"`
	Used       float64 `json:"Used"`
	Percentage float64 `json:"Percentage"`
}

type DiskInfo struct {
	Metric
	Drive      string `json:"drive"`
	TotalSize  uint64 `json:"totalSize"`
	Free       uint64 `json:"free"`
	Used       uint64 `json:"used"`
	FormatSize string `json:"format_size"`
	FormatFree string `json:"format_free"`
}

type SystemInfo struct {
	CPU    []CPUInfo    `json:"cpu"`
	Memory []MemoryInfo `json:"memory"`
	Disk   []DiskInfo   `json:"disk"`
}

type HealthStatus struct {
	SystemInfo SystemInfo `json:"system_info"`
	Uptime     string     `json:"uptime"`
	AgentInfo  AgentInfo  `json:"agent_info"`
}

// AgentInfo represents the complete metrics information for an agent
type AgentInfo struct {
	Version    string     `json:"version"`
	AgentID    string     `json:"agent_id"`
	IPAddress  string     `json:"IPAddress"`
	Name       string     `json:"name"`
	OS         string     `json:"os"`
	LastSync   *time.Time `json:"lastSync"`
	SDKVersion string     `json:"SDKVersion"`
}

// NetworkInfo holds network I/O statistics.
type NetworkInfo struct {
	Metric
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}
