package frontend

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type dashboard struct {
	app.Compo
	stats   Stats
	history []Stats
}

func (h *dashboard) Render() app.UI {
	return app.Div().Body(
		app.Div().
			Style("display", "flex").
			Body(
				&CPUGauge{Stats: h.stats},
				&MemoryGauge{Stats: h.stats},
				&CGoCalls{Stats: h.stats},
				&Goroutines{Stats: h.stats},
			),
	)
}

func (h *dashboard) OnMount(ctx app.Context) {
	ctx.Async(func() {
		for {
			// http request
			origin := app.Window().Get("location").Get("origin").String()
			resp, err := http.Get(origin + "/_goliveapi")
			if err != nil {
				panic(err)
			}
			// parse json
			dec := json.NewDecoder(resp.Body)
			stats := Stats{}
			if err := dec.Decode(&stats); err != nil {
				panic(err)
			}
			resp.Body.Close()
			ctx.Dispatch(func(ctx app.Context) {
				h.history = append(h.history, h.stats)
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
