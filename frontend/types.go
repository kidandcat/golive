package frontend

import (
	"runtime/debug"
	"runtime/metrics"
	"sync"
)

type Stats struct {
	sync.RWMutex `json:"-"`
	MemoryTotal  uint64
	MemoryUsed   uint64
	MemoryFree   uint64
	CPUTotal     uint64
	CPUUsed      uint64
	Uptime       uint64
	Metrics      []metrics.Sample
	GCstats      *debug.GCStats
	NumGoroutine int
	BuildInfo    *debug.BuildInfo
	Hostname     string
	NumCPU       int
	NumCgoCall   int64
	TraceData    []byte
}
