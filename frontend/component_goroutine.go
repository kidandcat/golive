package frontend

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Goroutine struct {
	app.Compo
	ID        int64
	State     string
	Reason    string
	Timestamp int64
}

func (h *Goroutine) Render() app.UI {
	elem := app.Div().Text(fmt.Sprintf("Goroutine %d: %s  (%s)", h.ID, h.State, h.Reason))
	if h.State == "Running" {
		elem = elem.Style("background-color", "green")
	}
	return elem
}
