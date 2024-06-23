package frontend

import (
	"math/rand"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Hello is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type Bar struct {
	app.Compo
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (h *Bar) Render() app.UI {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "My first bar chart",
		Subtitle: "It's extremely easy to use, right?",
	}))

	bar.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Category A", generateBarItems()).
		AddSeries("Category B", generateBarItems())

	snippets := bar.RenderSnippet()
	script := strings.Replace(snippets.Script, `<script type="text/javascript">`, ``, 1)
	script = strings.Replace(script, `</script>`, ``, 1)
	go func() {
		app.Window().Call("eval", script)
	}()
	return app.Div().Body(
		app.Raw(snippets.Element),
		app.Raw(snippets.Option),
	)
}

// generate random data for bar chart
func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}
