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
	router.Use(httplog.RequestLogger(app.Logger, []string{
		"/favicon.ico",
		"/status",
		"/ping",
	}))
	router.Use(middleware.Heartbeat("/ping"))

	router.NotFound(app.notFound)
	router.MethodNotAllowed(app.methodNotAllowed)

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	router.Handle("/static/*", fileServer)

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
		web.Get("/dashboard", app.dashboard)
		web.Post("/dashboard/clients", app.clients)
		web.Get("/dashboard/clients", app.clients)
		web.Get("/dashboard/clients/{id}", app.clientHome)
		web.Post("/dashboard/clients/{id}", app.clientHome)
		web.Get("/dashboard/clients/{id}/gitlab", app.clientGitlab)
		web.Post("/dashboard/clients/{id}/gitlab", app.clientGitlab)
		web.Get("/dashboard/clients/{id}/notifications", app.clientNotifications)
		web.Post("/dashboard/clients/{id}/notifications", app.clientNotifications)
	})
	router.Group(func(api chi.Router) {
		api.Route("/slash", func(slash chi.Router) {
			// slash specific api routes
		})
		// non-slash api routes
	})

	return router
}
