package golive

import (
	"encoding/json"
	"fmt"
	"golive/frontend"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/klauspost/compress/gzhttp"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const prefix = "golive"

var startTime time.Time

func init() {
	startTime = time.Now()
	frontend.Initialize()
	http.HandleFunc("/_goliveapi", handler)

	http.Handle(fmt.Sprintf("/%s/*", prefix), http.StripPrefix(fmt.Sprintf("/%s", prefix), gzhttp.GzipHandler(&app.Handler{
		Name:        "Go Live",
		Description: "Go monitoring dashboard",
		Resources:   ResourceFS(prefix),
		Scripts: []string{
			"/web/echarts.min.js",
		},
		Styles: []string{
			"/web/dark.css",
			"/web/style.css",
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

	memstats := new(runtime.MemStats)
	runtime.ReadMemStats(memstats)

	info, _ := debug.ReadBuildInfo()
	hostname, _ := os.Hostname()

	stats := frontend.Stats{
		MemoryTotal:  m.Total,
		MemoryUsed:   m.Used,
		MemoryFree:   m.Free,
		CPUTotal:     uint64(runtime.NumCPU() * 100),
		CPUUsed:      c.User + c.System,
		Uptime:       uint64(time.Since(startTime).Seconds()),
		Memstats:     memstats,
		NumGoroutine: runtime.NumGoroutine(),
		BuildInfo:    info,
		Hostname:     hostname,
		NumCPU:       runtime.NumCPU(),
		NumCgoCall:   runtime.NumCgoCall(),
	}

	jsonstats, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonstats)
}
