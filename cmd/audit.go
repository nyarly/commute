package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// cleanupCmd represents the cleanup command
	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "review state of the system vs. shared config",
		Long: longUsage(`Primarily, this prints out directories that seem to refer to a different remote
    than is recorded in the config. Generally, it's best to align remotes.  `),
		RunE:    auditFn,
	}
)

func init() {
	rootCmd.AddCommand(auditCmd)
}

func auditFn(cmd *cobra.Command, args []string) error {
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
      normal("%q: commute has %q but git remote is %q", path, rem, picked)
    }
  }
  return nil
}
