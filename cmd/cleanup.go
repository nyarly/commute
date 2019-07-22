package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	// cleanupCmd represents the cleanup command
	cleanupCmd = &cobra.Command{
		Use:   "cleanup",
		Short: "clean stale links to missing workspaces",
		Long: longUsage(`Over time, repository records get stale. For instance, if a workspace is deleted or moved,
		commute's idea of where it is will become wrong. This command checks that those records are current,
		and removes entries that are wrong.`),
		RunE:    cleanupFn,
	}
)

func init() {
	rootCmd.AddCommand(cleanupCmd)
  queryCommand(cleanupCmd)
}

func cleanupFn(cmd *cobra.Command, args []string) error {
	return filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if info.Mode()&os.ModeSymlink == 0 {
			return nil
		}
		rel, err := filepath.Rel(configDir, path)
		if err != nil {
			rel = path
		}

		rem := cfg.Remotes.findByLinkPath(path)
		if rem == remote("") {
			normal("removing %s: untracked\n", rel)
			return actually(func() error {
				return os.Remove(path)
			})
		}

		target, err := os.Readlink(path)
		if err != nil {
			verbose("%s link broken: %v\n", rel, err)
			return nil
		}

		if _, err := os.Stat(target); err != nil {
			normal("removing %s: target broken\n", rel)
			verbose("  (%s error: %v)\n", target, err)
			return actually(func() error {
				return os.Remove(path)
			})
		}

		verbose("  ok: %s -> %s\n", rel, target)

		return nil
	})
}
