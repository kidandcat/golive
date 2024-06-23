package frontend

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Hello is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type hello struct {
	app.Compo
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Line example"),
		&Line{},
		app.H1().Text("Bar example"),
		&Bar{},
	)
}
