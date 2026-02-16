package cmd

import (
	"fmt"
	"os"

	"terminal-brain/internal/app"

	"github.com/spf13/cobra"
)

var dbPath string

var rootCmd = &cobra.Command{
	Use:   "terminal-brain",
	Short: "Terminal Brain - a local-first cognitive workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dbPath == "" {
			return fmt.Errorf("db path cannot be empty")
		}
		return app.Run(dbPath)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&dbPath, "db", "terminal-brain.db", "path to sqlite database")
}
