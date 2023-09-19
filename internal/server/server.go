package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/danielmichaels/conrad/internal/config"
	"github.com/danielmichaels/conrad/internal/repository"
	"github.com/danielmichaels/conrad/internal/smtp"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Application struct {
	Config   *config.Conf
	Logger   zerolog.Logger
	wg       sync.WaitGroup
	Mailer   *smtp.Mailer
	Db       *repository.Queries
	Sessions *sessions.CookieStore
}

// Ptr takes in non-pointer and returns a pointer
func Ptr[T any](v T) *T {
	return &v
}

func (app *Application) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Server.Port),
		Handler:      app.routes(),
		IdleTimeout:  app.Config.Server.TimeoutIdle,
		ReadTimeout:  app.Config.Server.TimeoutRead,
		WriteTimeout: app.Config.Server.TimeoutWrite,
	}

	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.Logger.Warn().Str("signal", s.String()).Msg("caught signal")

		// Allow processes to finish with a ten-second window
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		app.Logger.Warn().Str("tasks", srv.Addr).Msg("completing background tasks")
		// Call wait so that the wait group can decrement to zero.
		app.wg.Wait()
		shutdownError <- nil
	}()
	app.Logger.Info().Str("server", srv.Addr).Msg("starting server")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		app.Logger.Warn().Str("server", srv.Addr).Msg("stopped server")
		return err
	}
	return nil
}

func Hash(plaintextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func Matches(plaintextPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func formatURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return url
}

// isInsecure checks the value of a form.Insecure and returns a bool as a string
func isInsecure(v string) string {
	if v == "on" || v == "true" {
		return "true"
	}
	return "false"
}

func urlIDParam(w http.ResponseWriter, r *http.Request) (int64, error) {
	param := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, err
}
