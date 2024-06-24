package frontend

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type MemoryGauge struct {
	app.Compo
	chartComponent
	Xaxis []string
	Stats Stats
}

func (h *MemoryGauge) Render() app.UI {
	gauge := charts.NewGauge()
	total := h.Stats.MemoryTotal
	value := h.Stats.MemoryUsed
	// percentage
	if total > 0 {
		value = value * 100 / total
	}
	gauge.AddSeries("Memory", []opts.GaugeData{{Name: "Memory", Value: value}})
	return h.renderChart(gauge)
}
