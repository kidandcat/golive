package frontend

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type CGoCalls struct {
	app.Compo
	chartComponent
	Xaxis []string
	Stats Stats
	prev  int64
}

func (h *CGoCalls) Render() app.UI {
	gauge := charts.NewGauge()
	gauge.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  "dark",
			Height: "300px",
			Width:  "300px",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
	)
	value := h.Stats.NumCgoCall - h.prev
	gauge.AddSeries("CGo/s", []opts.GaugeData{{Name: "CGo/s", Value: value}})
	h.prev = h.Stats.NumCgoCall
	return h.renderChart(gauge)
}
