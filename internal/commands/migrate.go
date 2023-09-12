package commands

import (
	"context"
	"fmt"
	"github.com/danielmichaels/conrad/assets"
	"os"

	"github.com/danielmichaels/conrad/internal/cmdutils"
	"github.com/danielmichaels/conrad/internal/config"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

func MigrateCmd(ctx context.Context) *cobra.Command {
	cfg := config.AppConfig()

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run the migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := cmdutils.NewLogger("migrate", cfg)

			db, err := cmdutils.NewDatabasePool(ctx, cfg)
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to open database")
				return err
			}
			defer db.Close()

			_ = goose.SetDialect("sqlite3")
			migrations := "assets/migrations"
			if os.Getenv("DOCKER") != "" {
				fmt.Println("running embedded migrations")
				goose.SetBaseFS(assets.EmbeddedFiles)
				migrations = "migrations"
			}
			err = goose.Up(db, migrations)
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to run migrations")
				return err
			}
			fmt.Println("migrations complete")
			return nil
		},
	}
	return cmd
}
