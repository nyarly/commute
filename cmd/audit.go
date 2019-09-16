package cmd

import (
  "fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	// cleanupCmd represents the cleanup command
	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "review state of the system vs. shared config",
		Long: longUsage(`Primarily, this prints out directories that seem to refer to a different remote
    than is recorded in the config. Generally, it's best to align remotes.  `),
		RunE: auditFn,
	}
)

func init() {
  auditCmd.Flags().Bool("header", true, "print headers for audit table")
	rootCmd.AddCommand(auditCmd)
}

func auditFn(cmd *cobra.Command, args []string) error {
	tabs := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
  if header, err := cmd.Flags().GetBool("header"); err == nil && header {
    fmt.Fprintln(tabs, "PATH\tCONFIGURED\tGIT REMOTE")
  }
	for _, rem := range cfg.Remotes {
		path, err := rem.localPath()
		if err != nil {
			continue
		}

		remotes, err := getRemotes(path)
		if err != nil {
			return err
		}

		picked, same := chooseRemote(cfg, remotes)
		if !same {
			fmt.Fprintf(tabs, "%s\t%s\t%s\n", path, rem, picked)
		}
	}
  tabs.Flush()
	return nil
}
