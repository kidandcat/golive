//go:build !wasm

package golive

import (
	"embed"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

//go:embed web
var web embed.FS

var _ app.ResourceResolver = (*embeddedResourceResolver)(nil)

func ResourceFS(prefix string) app.ResourceResolver {
	return embeddedResourceResolver{
		prefix:  prefix,
		Handler: http.FileServer(http.FS(web)),
	}
}

type embeddedResourceResolver struct {
	prefix string
	http.Handler
}

func (r embeddedResourceResolver) Resolve(location string) string {
	result := location
	if location == "" {
		result = "/" + r.prefix
	}
	if location[0] == '/' {
		result = "/" + r.prefix + location
	}
	return result
}
