package golive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/kidandcat/golive/frontend"
	"github.com/klauspost/compress/gzhttp"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"golang.org/x/exp/trace"
)

const prefix = "golive"
const BUFFER_SIZE = 1024 * 1024 * 10

var startTime time.Time

func init() {
	startTime = time.Now()
	frontend.Initialize()
	http.HandleFunc("/_goliveapi", handler)
	go collector()
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

var stats = frontend.Stats{}

func collector() {
	// Set up the flight recorder.
	fr := trace.NewFlightRecorder()
	fr.Start()
	for {
		c, err := cpu.Get()
		if err != nil {
			fmt.Printf("failed to get cpu stats: %v", err)
		}
		m, err := memory.Get()
		if err != nil {
			fmt.Printf("failed to get memory stats: %v", err)
		}

		memstats := new(runtime.MemStats)
		runtime.ReadMemStats(memstats)

		info, _ := debug.ReadBuildInfo()
		hostname, _ := os.Hostname()

		// https://pkg.go.dev/runtime/metrics#example-Read-ReadingAllMetrics

		// https://github.com/felixge/fgprof/blob/master/fgprof.go#L87

		// https://www.datadoghq.com/blog/go-memory-metrics/#how-to-analyze-go-memory-usage

		// Grab the snapshot.
		var b bytes.Buffer
		_, err = fr.WriteTo(&b)
		if err != nil {
			fmt.Printf("failed to write trace data: %v", err)
		}

		stats.Lock()
		stats.MemoryTotal = m.Total
		stats.MemoryUsed = m.Used
		stats.MemoryFree = m.Free
		stats.CPUTotal = uint64(runtime.NumCPU() * 100)
		stats.CPUUsed = c.User + c.System
		stats.Uptime = uint64(time.Since(startTime).Seconds())
		stats.Memstats = memstats
		stats.NumGoroutine = runtime.NumGoroutine()
		stats.BuildInfo = info
		stats.Hostname = hostname
		stats.NumCPU = runtime.NumCPU()
		stats.NumCgoCall = runtime.NumCgoCall()
		stats.TraceData = b.Bytes()
		stats.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	stats.RLock()
	jsonstats, err := json.Marshal(stats)
	stats.RUnlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonstats)
}
