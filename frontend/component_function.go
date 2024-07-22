package frontend

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Function struct {
	Stack   string
	Time    time.Duration
	Started map[int64]time.Time
}

func (h *Function) Render() app.UI {
	return app.Div().
		Style("border", "1px solid red").
		Style("margin", "5px").
		Text(fmt.Sprintf("Function [%s] %s", h.Time, h.Stack))
}
