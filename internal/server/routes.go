package server

import (
	"github.com/danielmichaels/conrad/assets"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"net/http"
)

func (app *Application) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(app.recoverPanic)
	router.Use(app.securityHeaders)
	router.Use(middleware.RealIP)
	router.Use(middleware.Compress(5))
	router.Use(httplog.RequestLogger(app.Logger))
	//router.Use(app.preventCSRF)
	//router.Use(app.authenticate)

	router.NotFound(app.notFound)
	router.MethodNotAllowed(app.methodNotAllowed)

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	router.Handle("/static/*", fileServer)

	router.Group(func(hook chi.Router) {
		hook.Post("/webhooks/gitlab", app.gitlabWebhook)
	})
	router.Get("/status", app.status)
	router.Group(func(web chi.Router) {
		web.Use(app.preventCSRF)
		web.Use(app.authenticate)
		web.Use(app.requireAnonymousUser)
		web.Get("/", app.home)
		web.Get("/signup", app.userSignup)
		web.Post("/signup", app.userSignup)
		web.Get("/login", app.userLogin)
		web.Post("/login", app.userLogin)
	})
	router.Group(func(web chi.Router) {
		web.Use(app.preventCSRF)
		web.Use(app.authenticate)
		web.Use(app.requireAuthenticatedUser)
		web.Post("/logout", app.userLogout)
		web.Get("/dashboard", app.home)
	})
	router.Group(func(api chi.Router) {
		// api routes
	})

	return router
}
