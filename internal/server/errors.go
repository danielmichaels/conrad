package server

import (
	"fmt"
	"github.com/danielmichaels/conrad/internal/render"
	"github.com/danielmichaels/conrad/internal/validator"
	"net/http"
)

// Error defines model for Error.
type Error struct {
	Body *map[string]interface{} `json:"body,omitempty"`
	// Code Error code
	Code int32 `json:"code"`
	// Status Error message
	Status string `json:"status"`
}

func (app *Application) errorMessage(w http.ResponseWriter, r *http.Request, e Error, headers http.Header) {
	err := render.JSONWithHeaders(w, int(e.Code), e, headers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	err = render.Page(w, http.StatusUnprocessableEntity, nil, "pages/500.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}

}

func (app *Application) notFound(w http.ResponseWriter, r *http.Request) {
	err := render.Page(w, http.StatusUnprocessableEntity, nil, "pages/404.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorMessage(w, r, Error{
		Code:   http.StatusMethodNotAllowed,
		Status: message,
	}, nil)
}

func (app *Application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorMessage(w, r, Error{
		Code:   http.StatusBadRequest,
		Status: err.Error(),
	}, nil)
}

func (app *Application) failedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator) {
	err := render.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) invalidAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", "Bearer")

	app.errorMessage(w, r, Error{
		Code:   http.StatusUnauthorized,
		Status: "Invalid authentication token",
	}, headers)
}

func (app *Application) authenticationRequired(w http.ResponseWriter, r *http.Request) {
	app.errorMessage(w, r, Error{
		Code:   http.StatusUnauthorized,
		Status: "You must be authenticated to access this resource",
	}, nil)
}

func (app *Application) basicAuthenticationRequired(w http.ResponseWriter, r *http.Request) {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	message := "You must be authenticated to access this resource"
	app.errorMessage(w, r, Error{
		Code:   http.StatusUnauthorized,
		Status: message,
	}, headers)
}

func (app *Application) errorHookEventNotFound(w http.ResponseWriter, r *http.Request) {
	app.errorMessage(w, r, Error{
		Code:   http.StatusBadRequest,
		Status: "webhook event not found",
	}, nil)
}
