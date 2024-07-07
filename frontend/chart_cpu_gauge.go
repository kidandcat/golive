package frontend

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type CPUGauge struct {
	app.Compo
	chartComponent
	Xaxis []string
	Stats Stats
	prev  uint64
}

func (h *CPUGauge) Render() app.UI {
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
	value := h.Stats.CPUUsed - h.prev
	h.prev = h.Stats.CPUUsed
	total := h.Stats.CPUTotal
	// percentage
	if total > 0 {
		value = value * 100 / total
	}
	gauge.AddSeries("CPU", []opts.GaugeData{{Name: "CPU", Value: value}})
	return h.renderChart(gauge)
}
