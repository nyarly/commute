package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("remote", "r", "", "specify a remote to add, rather than guess")
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

	var rem remote
	var found bool

	if cmd.Flags().Changed("remote") {
		chosenRemote, err := cmd.Flags().GetString("remote")
		if err != nil {
			return err
		}
		gr, err := pickNamedRemote(chosenRemote, remotes)
		if err != nil {
			return err
		}
		rem, found = chooseRemote(cfg, []gitRemote{gr})

	} else {
		rem, found = chooseRemote(cfg, remotes)
	}

	if !found {
		verbose("New remote %q added to configuration\n", rem)
		cfg.Remotes = append(cfg.Remotes, rem)
	}

	lp, err := rem.linkPath()
	if err != nil {
		return err
	}

	_, err = os.Stat(lp)
	if err == nil {
		p, _ := rem.localPath()
		return fmt.Errorf("remote already accounted for as %s", p)
	}
	os.Mkdir(filepath.Dir(lp), os.ModeDir|os.ModePerm)
	return os.Symlink(root, lp)
}

func pickNamedRemote(chosenRemote string, remotes []gitRemote) (gitRemote, error) {
	for _, r := range remotes {
		if chosenRemote == r.name {
			return r, nil
		}
	}

	remnames := []string{}
	for _, r := range remotes {
		remnames = append(remnames, r.name)
	}

	return gitRemote{}, fmt.Errorf("no remote named %q in local repo - found: %s", chosenRemote, strings.Join(remnames, ", "))
}
