package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add the current directory to the commute list",
	Long: longUsage(`Records the origin or upstream remote of the current git workspace
		in the commute list (if not already present) and marks the current directory
		as the canonical location of that repo locally.`),
	RunE: addFn,
}

func addFn(cmd *cobra.Command, args []string) error {
	root, err := getRepoRoot()
	if err != nil {
		return err
	}

	remotes, err := getRemotes(root)
	if err != nil {
		return err
	}

	rem, found := chooseRemote(cfg, remotes)

	if !found {
		cfg.Remotes = append(cfg.Remotes, rem)
		if e := cfg.save(); e != nil {
			return e
		}
	}

	_, err = os.Stat(rem.linkPath())
	if err == nil {
		p, _ := rem.localPath()
		return fmt.Errorf("remote already accounted for as %s", p)
	}
	os.Mkdir(filepath.Dir(rem.linkPath()), os.ModeDir|os.ModePerm)
	return os.Symlink(root, rem.linkPath())
}
