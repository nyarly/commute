package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func doAdd() error {
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
