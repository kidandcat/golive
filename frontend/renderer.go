package frontend

import (
	"regexp"
	"strings"

	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Snippet interface {
	RenderSnippet() render.ChartSnippet
}

type chartComponent struct {
	id      string
	element string
	options string
}

func (cc *chartComponent) renderChart(c Snippet) app.UI {
	snippets := c.RenderSnippet()
	script := strings.Replace(snippets.Script, `<script type="text/javascript">`, `try{`, 1)
	script = strings.Replace(script, `</script>`, `}catch(e){console.log(e)}`, 1)
	id := regexp.MustCompile(`id="([a-zA-Z0-9]+)"`).FindStringSubmatch(snippets.Element)[1]
	if cc.id != "" {
		script = strings.ReplaceAll(script, id, cc.id)
	} else {
		cc.id = id
		cc.element = snippets.Element
		cc.options = snippets.Option
	}
	go func() {
		app.Window().Call("eval", script)
	}()
	return app.Div().Body(
		app.Raw(cc.element),
	)
}
