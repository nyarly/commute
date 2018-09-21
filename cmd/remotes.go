package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

type (
	remote string

	remotes []remote
)

func (r *remote) name() (string, error) {
	m := remoteNameRE.FindStringSubmatch(string(*r))
	if m == nil || len(m) < 2 {
		return "", fmt.Errorf("badly formatted git remote: %s", *r)
	}
	return m[1], nil
}

func (r *remote) linkPath() (string, error) {
	name, err := r.name()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, name), nil
}

func (r *remote) localPath() (string, error) {
	path, err := r.linkPath()
	if err != nil {
		return "", err
	}

	return os.Readlink(path)
}

func (rs remotes) findByLinkPath(p string) remote {
	for _, r := range rs {
		if l, err := r.linkPath(); err == nil && l == p {
			return r
		}
	}

	return remote("")
}
