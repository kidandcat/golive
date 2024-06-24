package golive

import (
	"encoding/json"
	"golive/frontend"
	"net/http"
	"runtime"
	"time"

	"github.com/klauspost/compress/gzhttp"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

var startTime time.Time

func init() {
	startTime = time.Now()
	frontend.Initialize()
	http.HandleFunc("/_goliveapi", handler)

	http.Handle("/_golive/*", http.StripPrefix("/_golive", gzhttp.GzipHandler(&app.Handler{
		Name:        "Go Live",
		Description: "Go monitoring dashboard",
		Resources:   ResourceFS("_golive"),
		Scripts: []string{
			"https://go-echarts.github.io/go-echarts-assets/assets/echarts.min.js",
		},
	})))
}

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := cpu.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m, err := memory.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// d, err := disk.Get()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	stats := frontend.Stats{
		MemoryTotal: m.Total,
		MemoryUsed:  m.Used,
		MemoryFree:  m.Free,
		CPUTotal:    uint64(runtime.NumCPU() * 100),
		CPUUsed:     c.User + c.System,
		Uptime:      uint64(time.Since(startTime).Seconds()),
	}

	jsonstats, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonstats)
}
