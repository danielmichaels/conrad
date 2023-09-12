package server

import (
	"context"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// If there's a panic, set "Connection: close" on the response.
				// This will tell Go's HTTP server to automatically close the
				// current connection after the response has been sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has a type interface{}, so we
				// use fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. This will log the error using
				// our custom Logger type at the ERROR level and send the client
				// a 500 status.
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func (app *Application) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.Sessions.Get(r, "session")
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		userID, ok := session.Values["userID"].(int64)
		if ok {
			user, err := app.Db.GetUserById(context.Background(), userID)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			if user.ID != 0 {
				r = contextSetAuthenticatedUser(r, &user)
			}
		}

		next.ServeHTTP(w, r)
	})
}
func (app *Application) preventCSRF(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		MaxAge:   86400,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	})
	return csrfHandler
}
func (app *Application) requireAnonymousUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser := contextGetAuthenticatedUser(r)
		if authenticatedUser != nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Application) requireBasicAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, plaintextPassword, ok := r.BasicAuth()
		if !ok {
			app.basicAuthenticationRequired(w, r)
			return
		}

		//if app.config.basicAuth.username != username {
		if "admin" != username {
			app.basicAuthenticationRequired(w, r)
			return
		}

		if "password" != plaintextPassword {
			app.basicAuthenticationRequired(w, r)
			return
		}
		//err := bcrypt.CompareHashAndPassword([]byte(app.config.basicAuth.hashedPassword), []byte(plaintextPassword))
		//switch {
		//case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		//	app.basicAuthenticationRequired(w, r)
		//	return
		//case err != nil:
		//	app.serverError(w, r, err)
		//	return
		//}
		next.ServeHTTP(w, r)
	})
}
func (app *Application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser := contextGetAuthenticatedUser(r)
		if authenticatedUser == nil {
			session, err := app.Sessions.Get(r, "session")
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			session.Values["redirectPathAfterLogin"] = r.URL.Path
			err = session.Save(r, w)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
