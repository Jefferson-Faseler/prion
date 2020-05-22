package bundle

import (
	"io/ioutil"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
)

// DirPath is the location of the bundle directory
func DirPath() string {
	return viper.GetString("VIM_BUNDLE_DIR")
}

// Clone will download a copy of the git repository to the directory specified
func Clone(pkgURL, dirPath string) error {
	_, err := git.PlainClone(dirPath, false, &git.CloneOptions{
		URL:      pkgURL,
		Progress: os.Stdout,
	})
	return err
}

// Pull the latest changes from the remote repository of the directory given
func Pull(dirPath string) error {
	repo, err := git.PlainOpen(dirPath)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})

	return err
}

// Packages will return all the packages in the current user's bundle directory
func Packages() (pkgs []string, err error) {
	dirs, err := ioutil.ReadDir(viper.GetString("VIM_BUNDLE_DIR"))
	if err != nil {
		return pkgs, err
	}

	for _, dir := range dirs {
		pkgs = append(pkgs, dir.Name())
	}
	return pkgs, err
}
