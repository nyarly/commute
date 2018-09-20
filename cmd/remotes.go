package cmd

import (
	"os"
	"path/filepath"
)

type (
	remote string

	remotes []remote
)

func (r *remote) name() string {
	m := remoteNameRE.FindStringSubmatch(string(*r))
	if m == nil || len(m) < 2 {
		panic("Badly formatted git remote: " + string(*r))
	}
	return m[1]
}

func (r *remote) linkPath() string {
	return filepath.Join(configDir, r.name())
}

func (r *remote) localPath() (string, error) {
	return os.Readlink(r.linkPath())
}

func (rs remotes) findByLinkPath(p string) remote {
	for _, r := range rs {
		if r.linkPath() == p {
			return r
		}
	}

	return remote("")
}
