package frontend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type dashboard struct {
	app.Compo
	stats Stats
}

func (h *dashboard) Render() app.UI {
	return &Trace{Stats: h.stats}
}

func (h *dashboard) OnMount(ctx app.Context) {
	h.stats = Stats{}
	ctx.Async(func() {
		for {
			// http request
			origin := app.Window().Get("location").Get("origin").String()
			resp, err := http.Get(origin + "/_goliveapi")
			if err != nil {
				fmt.Printf("failed to get stats: %v", err)
				time.Sleep(2 * time.Second)
				ctx.Reload()
			}
			// parse json
			dec := json.NewDecoder(resp.Body)
			stats := Stats{}
			if err := dec.Decode(&stats); err != nil {
				fmt.Printf("failed to decode stats: %v", err)
				time.Sleep(2 * time.Second)
				ctx.Reload()
			}
			resp.Body.Close()
			ctx.Dispatch(func(ctx app.Context) {
				h.stats = stats
			})
			time.Sleep(1 * time.Second)
		}
	})
}

func (h *dashboard) OnAppUpdate(ctx app.Context) {
	if ctx.AppUpdateAvailable() {
		ctx.Reload()
	}
}
