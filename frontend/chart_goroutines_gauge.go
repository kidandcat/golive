package frontend

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Goroutines struct {
	app.Compo
	chartComponent
	Xaxis []string
	Stats Stats
}

func (h *Goroutines) Render() app.UI {
	gauge := charts.NewGauge()
	gauge.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  "dark",
			Height: "250px",
			Width:  "250px",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
	)
	gauge.AddSeries("Goroutines", []opts.GaugeData{{Name: "Goroutines", Value: h.Stats.NumGoroutine}})
	return h.renderChart(gauge)
}
