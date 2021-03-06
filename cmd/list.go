package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
  listCmd.Flags().BoolP("include-remotes", "r", false, "list remotes as well as directories")
	rootCmd.AddCommand(listCmd)
  queryCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list known repos and their status",
	Long: longUsage(`For every known repo, either prints the local workspace path,
		or an "missing" error to stderr.`),
	RunE: listFn,
}

func listFn(cmd *cobra.Command, args []string) error {
	for _, remote := range cfg.Remotes {
		lp, err := remote.linkPath()
		if err != nil {
			return err
		}

		_, err = os.Stat(lp)
		if err != nil {
			normal("%s -> MISSING\n", remote)
			continue
		}
		p, err := remote.localPath()
		if err != nil {
			normal("%s : %s\n", remote, err)
		}
    if include, err := cmd.Flags().GetBool("include-remotes"); err == nil && include {
      fmt.Printf("%s -> %s\n", remote, p)
    } else {
      fmt.Printf("%s\n", p)
    }
	}
	return nil
}
