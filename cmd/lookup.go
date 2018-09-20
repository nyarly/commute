package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:   "lookup <remote>",
	Short: "find the local workspace directory for a remote",
	Long:  longUsage(``),
	Args:  cobra.ExactArgs(1),
	RunE:  lookupFn,
}

func init() {
	rootCmd.AddCommand(lookupCmd)
}

func lookupFn(cmd *cobra.Command, args []string) error {
	rem := remote(args[0])

	path, err := rem.localPath()
	if err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}
