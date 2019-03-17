package wyag

import (
	"fmt"
	"os"
	pathpkg "path"

	"github.com/go-ini/ini"
)

type GitRepository struct {
	workTree string
	gitdir   string
	config   *ini.File
}

func NewGitRepository(path string, force bool) (*GitRepository, error) {
	gr := new(GitRepository)

	gr.workTree = path
	gr.gitdir = pathpkg.Join(path, ".git")

	gd, err := os.Lstat(gr.gitdir)
	if err != nil {
		return nil, err
	}

	if !force && !gd.IsDir() {
		return nil, fmt.Errorf("Not a Git repository %s", path)
	}

	// cfg, err := ini.Load("config")
	cfPath, err := repoFile(gr, "config", false)
	if err != nil {
		return nil, err
	}

	if cfPath != "" && Exists(cfPath) {
		gr.config, err = ini.Load(cfPath)
	} else if !force {
		return nil, fmt.Errorf("Configuration file missing")
	}

	if !force {
		vers, err := gr.config.Section("core").Key("repositoryformatversion").Int()
		if err != nil {
			return nil, err
		}

		if (vers != 0) && !force {
			return nil, fmt.Errorf("Unsupported repositoryformatversion %d", vers)
		}
	}

	return gr, nil
}

// Compute path under repo's gitdir.
func repoPath(gr *GitRepository, path string) string {
	return pathpkg.Join(gr.gitdir, path)
}

func repoFile(gr *GitRepository, path string, mkDir bool) (string, error) {
	_, err := repoDir(gr, path, mkDir)
	if err != nil {
		return "", err
	}
	return repoPath(gr, path), nil
}

func repoDir(gr *GitRepository, path string, mkDir bool) (string, error) {
	p := repoPath(gr, path)

	if Exists(p) {
		if IsDir(p) {
			return p, nil
		} else {
			return "", fmt.Errorf("Not a directory %s", p)
		}
	}

	if mkDir {
		err := os.Mkdir(p, 0777)
		if err != nil {
			return "", err
		}
		return p, nil
	} else {
		return "", nil
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil || !f.IsDir() {
		return false
	}
	return true
}
