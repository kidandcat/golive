package frontend

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func Initialize() {
	app.Route("/", func() app.Composer { return &dashboard{} })
	app.Route("/hello", func() app.Composer { return &dashboard{} })
	app.RunWhenOnBrowser()
}
