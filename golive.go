package golive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/metrics"
	"time"

	"github.com/kidandcat/golive/frontend"
	"github.com/klauspost/compress/gzhttp"
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
	descs := metrics.All()
	samples := make([]metrics.Sample, len(descs))
	for i := range samples {
		samples[i].Name = descs[i].Name
	}
	// Set up the flight recorder.
	fr := trace.NewFlightRecorder()
	fr.SetPeriod(999 * time.Millisecond)
	fr.Start()
	var b bytes.Buffer
	var err error
	var info *debug.BuildInfo
	var ncpu int
	var hostname string
	var n int
	for {
		if n == math.MaxInt-1 {
			n = 0
		}

		if n%10 == 0 {
			info, _ = debug.ReadBuildInfo()
			ncpu = runtime.NumCPU()
			hostname, _ = os.Hostname()
		}

		metrics.Read(samples)
		// https://github.com/felixge/fgprof/blob/master/fgprof.go#L87
		// https://www.datadoghq.com/blog/go-memory-metrics/#how-to-analyze-go-memory-usage

		// Grab the snapshot.
		b.Reset()
		_, err = fr.WriteTo(&b)
		if err != nil {
			fmt.Printf("failed to write trace data: %v", err)
		}

		uptime := uint64(time.Since(startTime).Seconds())
		ngos := runtime.NumGoroutine()
		ccgo := runtime.NumCgoCall()
		stats.Lock()
		stats.Metrics = samples
		stats.Uptime = uptime
		stats.NumGoroutine = ngos
		stats.BuildInfo = info
		stats.Hostname = hostname
		stats.NumCPU = ncpu
		stats.NumCgoCall = ccgo
		stats.TraceData = b.Bytes()
		stats.Unlock()
		time.Sleep(999 * time.Millisecond)
		n++
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
