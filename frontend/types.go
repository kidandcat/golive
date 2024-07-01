package frontend

import (
	"runtime"
	"runtime/debug"
)

type Stats struct {
	MemoryTotal  uint64
	MemoryUsed   uint64
	MemoryFree   uint64
	CPUTotal     uint64
	CPUUsed      uint64
	Uptime       uint64
	Memstats     *runtime.MemStats
	GCstats      *debug.GCStats
	NumGoroutine int
	BuildInfo    *debug.BuildInfo
	Hostname     string
	NumCPU       int
	NumCgoCall   int64
}
