package frontend

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Trace struct {
	app.Compo
	Xaxis []string
	Stats Stats
}

func (h *Trace) Render() app.UI {
	return app.Div().ID("trace").Text(fmt.Sprintf("%+v", h.Stats.TraceEvents))
}
