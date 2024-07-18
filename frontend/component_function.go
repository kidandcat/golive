package frontend

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Function struct {
	app.Compo
	Name           string
	State          string
	AccumulateTime float64
	Timestamp      int64
}

func (h *Function) Render() app.UI {
	elem := app.Div().Text(fmt.Sprintf("Function %s: %s", h.Name, h.State))
	if h.State == "Running" {
		elem = elem.Style("background-color", "green")
	}
	return elem
}
