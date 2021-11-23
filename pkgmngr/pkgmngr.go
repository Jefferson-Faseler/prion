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
	return install(pkgURL, reinstall)
}

func install(pkgURL string, reinstall bool) error {
	pkgName := getPkgName(pkgURL)
	dirPath := filepath.Join(bundle.DirPath(), pkgName)

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

// Reinstall will re-clone a package in place
func Reinstall(pkgName string) error {
	dirPath := filepath.Join(bundle.DirPath(), pkgName)

	remoteURL, err := bundle.RemoteURL(dirPath)
	if err != nil {
		return err
	}

	return install(remoteURL, true)
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
			fmt.Println(fmt.Sprintf(`%s is having troubles updating.
This can happen when the package was shallow installed, and needs repairing.
Try reinstalling the package with the url to see if this fixes the issue.

To reinstall simply run:
prion reinstall %s

`, pkgName, pkgName))
		}
		if strings.Contains(err.Error(), "ssh: handshake failed: knownhosts: key mismatch") {
			userHomeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Could not determine user's home directory")
				return false, err
			}
			fmt.Println(fmt.Sprintf(`%s cannot connect to its remote host over ssh.
If you have not set up ssh connections with your git server you will need to
follow the appropriate guide to do so. If you have set up ssh and this error
does not appear when you git clone using ssh then you likely have an outdated
or non-prioritized ssh fingerprint from the host that prion's dependencies are
attempting to use. This could be from the host updating its public key
fingerprint or using different key types; rsa or ed25519.


You can attempt to resolve this by running the below command to update your
known hosts for you. Replace github.com if you are using a different remote
host.

ssh-keyscan -H github.com >> %s/.ssh/known_hosts
`, pkgName, userHomeDir))
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
