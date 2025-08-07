package models

import "time"

type ProcInfo struct {
	PID        int32
	Name       string
	CPUPercent float64
	MemPercent float32
	User       string
	Status     string
	CreateTime time.Time
}
