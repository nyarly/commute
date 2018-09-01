package main

import (
	"fmt"
	"os"
)

func doList() error {
	for _, remote := range cfg.Remotes {
		_, err := os.Stat(remote.linkPath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s -> MISSING\n", remote)
			continue
		}
		p, err := remote.localPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s : %s\n", remote, err)
		}
		fmt.Printf("%s\n", p)
	}
	return nil
}
