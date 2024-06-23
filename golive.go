package golive

import (
	"golive/frontend"
	"net/http"

	"github.com/klauspost/compress/gzhttp"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func init() {
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
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello World _goliveapi"))
}
