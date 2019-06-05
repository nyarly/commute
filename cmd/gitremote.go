package cmd

import (
	"os"
	"os/exec"
	"strings"
)

type gitRemote struct {
	name, url string
}

func getRepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return lookup(wd, `.git`)
}

func chooseRemote(cfg *config, remotes []gitRemote) (remote, bool) {
	var best gitRemote

	for _, rr := range remotes {
    verbose("Considering remote: %s %s", rr.name, rr.url)
		for _, rem := range cfg.Remotes {
			if sameRemote(remote(rr.url), rem) {
        verbose("Remote exists in cfg: %q. Using that.", rem)
				return remote(rr.url), true
			}
		}
		if rr.name == `origin` ||
			(rr.name == `upstream` && best.name != `origin`) ||
			best.url == `` {
        verbose("We like %s better than %s", rr.url, best.url)
        best = rr
		}
	}

  verbose("Chose %s", best.url)

	return remote(best.url), false
}

func sameRemote(l, r remote) bool {
	return l == r || l+".git" == r || l == r+".git"
}

func getRemotes(root string) ([]gitRemote, error) {
	git := exec.Command(`git`, `remote`, `-v`)
	git.Dir = root

	rems := []gitRemote{}
	out, err := git.CombinedOutput()
	if err != nil {
		return rems, err
	}

	for _, line := range strings.Split(string(out), "\n") {
		fields := fieldsRE.Split(line, 3)
		if len(fields) < 2 || fields[2] != `(fetch)` {
			continue
		}
		rems = append(rems, gitRemote{name: fields[0], url: fields[1]})
	}
	return rems, nil
}
