package pkgmngr

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jefferson-Faseler/prion/internal/bundle"
)

// Install is the entrypoint for installing a git repository to the vim bundle
func Install(pkgURL string) error {
	pkgName := getPkgName(pkgURL)
	dirPath := filepath.Join(bundle.DirPath(), pkgName)

	if isDirPresent(dirPath) {
		return errors.New(pkgName + ` is already installed
To reinstall remove the package first and then install.
Or to simply update run:
prion update ` + pkgName)
	}
	tmpDirPath := filepath.Join(os.TempDir(), pkgName)
	err := os.MkdirAll(tmpDirPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = bundle.Clone(pkgURL, tmpDirPath)
	if err != nil {
		rmErr := os.RemoveAll(tmpDirPath)
		if rmErr != nil {
			return rmErr
		}
		return err
	}

	return os.Rename(tmpDirPath, dirPath)
}

// Remove will remove a package froom the vim bundle
func Remove(pkgName string) error {
	dirPath := filepath.Join(bundle.DirPath(), pkgName)
	if isDirMissing(dirPath) {
		return errors.New("No package named " + pkgName)
	}
	return os.RemoveAll(dirPath)
}

// Update updates a vim package, returns if there were changes and errors
func Update(pkgName string) (bool, error) {
	dirPath := filepath.Join(bundle.DirPath(), pkgName)

	err := bundle.Pull(dirPath)
	if err != nil {
		if strings.Contains(err.Error(), "already up-to-date") {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func getPkgName(pkgURL string) string {
	// pkgURL is assumed to be formatted as
	// git@github.com:owner/repo.git
	// or
	// https://github.com/owner/repo.git

	splitAddress := strings.Split(pkgURL, "/")
	gitPkgName := splitAddress[len(splitAddress)-1]
	pkgName := strings.Split(gitPkgName, ".")[0]
	return pkgName
}

func isDirPresent(dirPath string) bool {
	_, pathErr := os.Stat(dirPath)

	// type *PathError is able to confirm why the error occured
	return !os.IsNotExist(pathErr)
}

func isDirMissing(dirPath string) bool {
	_, pathErr := os.Stat(dirPath)

	// type *PathError is able to confirm why the error occured
	return os.IsNotExist(pathErr)
}
