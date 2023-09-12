package commands

import (
	"context"

	"github.com/spf13/cobra"
)

func Execute(ctx context.Context) int {
	rootCmd := &cobra.Command{
		Use:   "conrad",
		Short: "conrad is a CLI for things",
		Long:  "conrad is a CLI for things. It does stuff and things.",
	}

	rootCmd.AddCommand(ServeCmd(ctx))
	rootCmd.AddCommand(MigrateCmd(ctx))

	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}
