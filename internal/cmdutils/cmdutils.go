package cmdutils

import (
	"context"
	"database/sql"
	"time"

	"github.com/danielmichaels/conrad/internal/config"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

func NewLogger(name string, cfg *config.Conf) zerolog.Logger {
	logger := httplog.NewLogger(name, httplog.Options{
		JSON:     cfg.AppConf.LogJson,
		Concise:  cfg.AppConf.LogConcise,
		LogLevel: cfg.AppConf.LogLevel,
	})
	if cfg.AppConf.LogCaller {
		logger = logger.With().Caller().Logger()
	}
	return logger
}

func NewDatabasePool(ctx context.Context, cfg *config.Conf) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.Db.DbName)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
