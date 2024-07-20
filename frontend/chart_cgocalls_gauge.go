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
			Height: "250px",
			Width:  "250px",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
	)
	setMinMax := charts.WithSeriesOpts(func(s *charts.SingleSeries) {
		// Gauge
		//  Progress *opts.Progress `json:"progress,omitempty"`
		//  AxisTick *opts.AxisTick `json:"axisTick,omitempty"`
		//  Detail   *opts.Detail   `json:"detail,omitempty"`
		//  Title    *opts.Title    `json:"title,omitempty"`
		//  Min      int            `json:"min,omitempty"`
		//  Max      int            `json:"max,omitempty"`

		//  Large               types.Bool `json:"large,omitempty"`
		//  LargeThreshold      int        `json:"largeThreshold,omitempty"`
		//  HoverLayerThreshold int        `json:"hoverLayerThreshold,omitempty"`
		//  UseUTC              types.Bool `json:"useUTC,omitempty"`
		s.Min = 1
		s.Max = 12
	})
	value := h.Stats.NumCgoCall - h.prev
	gauge.AddSeries("CGo/s", []opts.GaugeData{{Name: "CGo/s", Value: value}}, setMinMax)
	h.prev = h.Stats.NumCgoCall
	return h.renderChart(gauge)
}
