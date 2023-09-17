package commands

import (
	"context"
	"github.com/danielmichaels/conrad/internal/cmdutils"
	"github.com/danielmichaels/conrad/internal/config"
	"github.com/danielmichaels/conrad/internal/repository"
	"github.com/danielmichaels/conrad/internal/server"
	"github.com/danielmichaels/conrad/internal/smtp"
	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
	"net/http"
)

func ServeCmd(ctx context.Context) *cobra.Command {
	cfg := config.AppConfig()
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "start the webserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := cmdutils.NewLogger("server", cfg)
			if cfg.AppConf.LogCaller {
				logger = logger.With().Caller().Logger()
			}

			dbc, err := cmdutils.NewDatabasePool(ctx, cfg)
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to open database")
			}

			defer dbc.Close()
			db := repository.New(dbc)

			mailer := smtp.NewMailer(
				cfg.Smtp.Host,
				cfg.Smtp.Port,
				cfg.Smtp.Username,
				cfg.Smtp.Password,
				cfg.Smtp.Sender,
			)

			keyPairs := [][]byte{[]byte("eda3fcy34wwevtz2vartsjm3h3qbqbqu"), nil}
			sessionStore := sessions.NewCookieStore(keyPairs...)
			sessionStore.Options = &sessions.Options{
				HttpOnly: true,
				MaxAge:   86400 * 7,
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
				Secure:   true,
			}

			app := &server.Application{
				Config:   cfg,
				Logger:   logger,
				Mailer:   mailer,
				Db:       db,
				Sessions: sessionStore,
			}
			err = app.Serve()
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to start server")
			}
			return nil
		},
	}
	return cmd
}
