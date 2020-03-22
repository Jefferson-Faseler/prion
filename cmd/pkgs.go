package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jefferson-Faseler/prion/internal/bundle"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	all bool // used for updating all packages
)

// installPkgCmd represents the command for installing vim packages
var installPkgCmd = &cobra.Command{
	Use:           "install [pkg url]",
	Short:         "easily install a vim package",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return nil
		}

		for _, pkgURL := range args {
			pkgName := getPkgName(pkgURL)
			dirPath := filepath.Join(bundleDir(), pkgName)

			if isDirPresent(dirPath) {
				return fmt.Errorf(pkgName + ` is already installed
To reinstall remove the package first and then install.
Or to simply update run:
prion update ` + pkgName)
			}
			tmpDirPath := filepath.Join(os.TempDir(), pkgName)
			err := os.MkdirAll(tmpDirPath, os.ModePerm)
			if err != nil {
				return err
			}

			fmt.Println("Installing " + pkgURL)
			err = bundle.Clone(pkgURL, tmpDirPath)
			if err != nil {
				rmErr := os.RemoveAll(tmpDirPath)
				if rmErr != nil {
					return rmErr
				}
				return err
			}
			err = os.Rename(tmpDirPath, dirPath)
			return err
		}
		return nil
	},
}

// removePkgCmd represents the command for removing vim packages
var removePkgCmd = &cobra.Command{
	Use:           "rm [pkg name]",
	Aliases:       []string{"remove"},
	Short:         "easily remove a vim package",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return nil
		}

		for _, pkgName := range args {
			dirPath := filepath.Join(bundleDir(), pkgName)

			if isDirMissing(dirPath) {
				return errors.New("No package named " + pkgName)
			}
			err := os.RemoveAll(dirPath)
			if err != nil {
				return err
			}

			fmt.Println(pkgName + " removed")
		}
		return nil
	},
}

// updatePkgCmd represents the command for updating vim packages
var updatePkgCmd = &cobra.Command{
	Use:           "update [pkg url]",
	Short:         "easily update a vim package",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var pkgs []string
		var err error

		if len(args) == 0 && all == false {
			cmd.Help()
			return nil
		} else if all == true {
			_pkgs, err := bundle.Packages()
			if err != nil {
				return err
			}
			pkgs = _pkgs
		} else {
			pkgs = args
		}

		for _, pkgName := range pkgs {
			dirPath := filepath.Join(bundleDir(), pkgName)

			// will return an error if the dir is missing
			err = bundle.Pull(dirPath)
			if err != nil {
				if strings.Contains(err.Error(), "already up-to-date") {
					fmt.Println(err)
				} else {
					return err
				}
			} else {
				fmt.Println("Updating " + pkgName)
			}
		}
		return nil
	},
}

// listPkgCmd represents the command for updating vim packages
var listPkgCmd = &cobra.Command{
	Use:           "ls [pkg url]",
	Aliases:       []string{"list"},
	Short:         "list all your vim packages",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgs, err := bundle.Packages()
		if err != nil {
			return err
		}

		for _, pkgName := range pkgs {
			fmt.Println(pkgName)
		}
		return nil
	},
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

func bundleDir() string {
	return viper.GetString("VIM_BUNDLE_DIR")
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

func init() {
	rootCmd.AddCommand(installPkgCmd)
	rootCmd.AddCommand(removePkgCmd)
	rootCmd.AddCommand(updatePkgCmd)
	rootCmd.AddCommand(listPkgCmd)

	updatePkgCmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "all packages")
}
