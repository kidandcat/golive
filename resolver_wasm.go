//go:build wasm

package golive

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func ResourceFS() app.ResourceResolver {
	return app.LocalDir("web")
}
