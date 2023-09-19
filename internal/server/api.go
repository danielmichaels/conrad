package server

import (
	"github.com/danielmichaels/conrad/internal/render"
	"github.com/danielmichaels/conrad/internal/version"
	"net/http"
)

func (app *Application) status(w http.ResponseWriter, r *http.Request) {
	data := Map{
		"status":  "OK",
		"version": version.Get(),
	}
	err := render.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
