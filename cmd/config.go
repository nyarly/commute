package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "subcommands for updating synchronized config",
	Long: longUsage(``),
  SilenceUsage: true,
}
