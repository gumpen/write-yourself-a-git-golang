package wyag

import (
	"fmt"
	"io/ioutil"
	"os"
	pathpkg "path"

	"github.com/go-ini/ini"
)

// GitRepository is the struct of git repository
type GitRepository struct {
	workTree string
	gitdir   string
	config   *ini.File
}

type GitObject struct {
	repository string
}

// NewGitRepository is constructor of GitRepository struct
func NewGitRepository(path string, force bool) (*GitRepository, error) {
	gr := new(GitRepository)

	gr.workTree = path
	gr.gitdir = pathpkg.Join(path, ".git")

	if !force && !IsDir(gr.gitdir) {
		return nil, fmt.Errorf("Not a Git repository %s", path)
	}

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

// NewGitObject is constructor of GitObject struct
func NewGitObject(repository string, data string) (*GitObject, error) {
	gob := new(GitObject)
	gob.repository = repository

	if data != "" {
		err := gob.deserialize(data)
		if err != nil {
			return nil, err
		}
	}

	return gob, nil
}

func (gob *GitObject) serialize(data string) error {
	return fmt.Errorf("Unimplemented")
}

func (gob *GitObject) deserialize(data string) error {
	return fmt.Errorf("Unimplemented")
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
		err := os.MkdirAll(p, 0777)
		if err != nil {
			return "", err
		}
		return p, nil
	} else {
		return "", nil
	}
}

func (gr *GitRepository) repoDefaultConfig(path string) error {
	_, err := gr.config.NewSection("core")
	if err != nil {
		return err
	}
	gr.config.Section("core").Key("repositoryformatversion").SetValue("0")
	gr.config.Section("core").Key("filemode").SetValue("false")
	gr.config.Section("core").Key("bare").SetValue("false")

	gr.config.SaveTo(path)

	return nil
}

// Create a new repository at path.
func repoCreate(path string) (*GitRepository, error) {
	gr, err := NewGitRepository(path, true)
	if err != nil {
		return nil, err
	}

	if Exists(gr.workTree) {
		if !IsDir(gr.workTree) {
			return nil, fmt.Errorf("%s is not a directory", path)
		}

		list, err := ioutil.ReadDir(gr.workTree)
		if err != nil {
			return nil, err
		}
		if len(list) != 0 {
			return nil, fmt.Errorf("%s is not empty", path)
		}
	} else {
		if err = os.Mkdir(gr.workTree, 0777); err != nil {
			return nil, err
		}
	}

	repoDir(gr, "branches", true)
	repoDir(gr, "objects", true)
	repoDir(gr, "refs/tags", true)
	repoDir(gr, "refs/heads", true)

	descriptionPath, err := repoFile(gr, "description", false)
	if err != nil {
		return nil, err
	}
	descriptionFile, err := os.Create(descriptionPath)
	if err != nil {
		return nil, err
	}
	defer descriptionFile.Close()
	descriptionFile.WriteString("Unnamed repository; edit this file 'description' to name the repository.\n")

	headPath, err := repoFile(gr, "HEAD", false)
	if err != nil {
		return nil, err
	}
	headFile, err := os.Create(headPath)
	if err != nil {
		return nil, err
	}
	defer headFile.Close()
	headFile.WriteString("ref: refs/heads/master\n")

	configPath, err := repoFile(gr, "config", false)
	if err != nil {
		return nil, err
	}
	configFile, err := os.Create(configPath)
	if err != nil {
		return nil, err
	}
	if err := configFile.Close(); err != nil {
		return nil, err
	}

	gr.config, err = ini.Load(configPath)
	if err != nil {
		return nil, err
	}
	err = gr.repoDefaultConfig(configPath)
	if err != nil {
		return nil, err
	}

	return gr, nil
}

func repoFind(path string, required bool) (*GitRepository, error) {
	if path == "" {
		path = "."
	}

	path, err := RealPath(path)
	if err != nil {
		return nil, err
	}

	if IsDir(pathpkg.Join(path, ".git")) {
		gr, err := NewGitRepository(path, false)
		if err != nil {
			return nil, err
		}
		return gr, nil
	}

	parent, err := RealPath(pathpkg.Join(path, ".."))
	if err != nil {
		return nil, err
	}

	if parent == path {
		if required {
			return nil, fmt.Errorf("No git directory")
		}
		return nil, nil
	}

	return repoFind(parent, required)
}

// Exists check
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir check
func IsDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil || !f.IsDir() {
		return false
	}
	return true
}

func RealPath(path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if fi.Mode() == os.ModeSymlink {
		originFilePath, err := os.Readlink(fi.Name())
		if err != nil {
			return "", err
		}
		return originFilePath, nil
	}
	return path, nil
}
