package cmd

import (
	"fmt"
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
		PreRunE: preCleanup,
		RunE:    cleanupFn,
	}

	quietFlag bool
	verbFlag  bool
	dryFlag   bool

	normal   = func(fmt string, as ...interface{}) {}
	verbose  = func(fmt string, as ...interface{}) {}
	actually = func(f func() error) error { return f() }
)

func init() {
	rootCmd.AddCommand(cleanupCmd)
	cleanupCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "quiet output")
	cleanupCmd.Flags().BoolVarP(&verbFlag, "verbose", "v", false, "verbose output")
	cleanupCmd.Flags().BoolVarP(&dryFlag, "dryrun", "d", false, "take no action")
}

func preCleanup(cmd *cobra.Command, args []string) error {
	if verbFlag && quietFlag {
		return fmt.Errorf("can't specify both quiet and verbose")
	}

	if !quietFlag {
		normal = func(tmpl string, as ...interface{}) {
			fmt.Printf(tmpl, as...)
		}
	}
	if verbFlag {
		verbose = func(tmpl string, as ...interface{}) {
			fmt.Printf(tmpl, as...)
		}
	}

	if dryFlag {
		actually = func(func() error) error { return nil }
	}
	return nil
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
