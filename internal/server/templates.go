package server

import (
	"github.com/danielmichaels/conrad/internal/version"
	"github.com/justinas/nosurf"
	"net/http"
	"time"
)

// Map is a helper for defining a map[string]any with brevity in mind.
type Map map[string]any

func (app *Application) newTemplateData(r *http.Request) map[string]any {
	data := map[string]any{
		"AppName":           "Conrad",
		"AuthenticatedUser": contextGetAuthenticatedUser(r),
		"CSRFToken":         nosurf.Token(r),
		"Version":           version.Get(),
		"CurrentYear":       time.Now().Year(),
	}
	return data
}

func (app *Application) newEmailData(r *http.Request) map[string]any {
	data := map[string]any{}

	return data
}
