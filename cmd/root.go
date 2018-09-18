package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute runs the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "commute",
	Short: "commute maintains the mapping of remote repositories to local workspaces",
	Long: longUsage(`A utility for keeping track of where repos are on you work machine,
		suitable for use by humans and scripts.`),
	PersistentPreRunE: setupStuff,
}

func setupStuff(cmd *cobra.Command, args []string) error {
	err := setupPaths()
	if err != nil {
		return err
	}

	return loadConfig()
}
