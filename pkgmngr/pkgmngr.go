package pkgmngr

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jefferson-Faseler/prion/internal/bundle"
)

// Install will clone the package URL to a temporary directory and move the contents to the vim bundle if everything was successful, otherwise it will attempt to delete the temp directory and return any errors
func Install(pkgURL string) error {
	var reinstall bool
	pkgName := getPkgName(pkgURL)
	dirPath := filepath.Join(bundle.DirPath(), pkgName)

	if isDirPresent(dirPath) {
		fmt.Printf(`%s is already installed, would you like to reinstall it? Enter 'y' to reinstall or any other key to cancel.
`, pkgName)
		reader := bufio.NewReader(os.Stdin)
		userInput, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if reinstall = strings.HasPrefix(userInput, "y"); reinstall != true {
			return nil
		}
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

	if reinstall {
		err = Remove(pkgName)
		if err != nil {
			return err
		}
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

// Update updates a vim package, returns a boolean denoting if the package was already up-to-date and the error if any occured
func Update(pkgName string) (bool, error) {
	dirPath := filepath.Join(bundle.DirPath(), pkgName)

	err := bundle.Pull(dirPath)
	if err != nil {
		if strings.Contains(err.Error(), "already up-to-date") {
			return true, nil
		}
		if strings.Contains(err.Error(), "object not found") {
			fmt.Println(pkgName + ` is having troubles updating.
This can happen when the package was shallow installed, and needs repairing.
Try reinstalling the package with the url to see if this fixes the issue.`)
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
	pkgName := strings.TrimRight(gitPkgName, ".git")
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
