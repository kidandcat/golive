package frontend

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Goroutine struct {
	app.Compo
	ID        string
	State     string
	Stack     string
	Reason    string
	Timestamp int64
	prevStack string
}

const (
	IgnoreGopark       = "runtime.gopark"
	IgnoreGoparkUnlock = "runtime.goparkunlock"
	IgnoreScavenger    = "runtime.(*scavengerState).park"
)

func (h *Goroutine) Render() app.UI {
	elem := app.Div().ID(h.ID)
	if h.Stack == "" {
		elem = elem.Style("color", "gray")
		h.Stack = h.prevStack
	} else {
		h.prevStack = h.Stack
	}
	elem.Text(fmt.Sprintf("Goroutine%s [%s %s]  %s", h.ID, h.State, h.Reason, h.Stack))
	if h.State == "Running" {
		elem = elem.Style("background-color", "green")
	} else if h.State != "Waiting" {
		elem = elem.Style("background-color", "purple")
	}
	return elem
}
