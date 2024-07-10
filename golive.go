package golive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/trace"
	"time"

	"github.com/kidandcat/golive/frontend"
	"github.com/klauspost/compress/gzhttp"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	gtrace "honnef.co/go/gotraceui/trace"
)

const prefix = "golive"
const BUFFER_SIZE = 1024 * 1024

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

var traceBuffer = new(bytes.Buffer)
var traceParser *gtrace.Parser
var stats = frontend.Stats{}

func collector() {
	var err error
	traceParser, err = gtrace.NewParser(traceBuffer)
	if err != nil {
		log.Fatalf("failed to create trace parser: %v", err)
	}
	if err := trace.Start(traceBuffer); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()
	for {
		c, err := cpu.Get()
		if err != nil {
			log.Printf("failed to get cpu stats: %v", err)
		}
		m, err := memory.Get()
		if err != nil {
			log.Printf("failed to get memory stats: %v", err)
		}

		memstats := new(runtime.MemStats)
		runtime.ReadMemStats(memstats)

		info, _ := debug.ReadBuildInfo()
		hostname, _ := os.Hostname()

		tt, err := traceParser.Parse()
		if err != nil {
			log.Printf("failed to parse trace: %v", err)
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
		if traceBuffer.Len() > BUFFER_SIZE {
			data := traceBuffer.Bytes()[traceBuffer.Len()-BUFFER_SIZE:]
			traceBuffer.Reset()
			traceBuffer = bytes.NewBuffer(data)
		}
		stats.TraceData = traceBuffer.Bytes()
		stats.TraceEvents = tt.Events
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
