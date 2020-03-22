package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jefferson-Faseler/prion/internal/bundle"
	"github.com/Jefferson-Faseler/prion/plugmngr"

	"github.com/spf13/cobra"
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
			err := plugmngr.Install(pkgURL)
			if err != nil {
				return err
			}
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
			dirPath := filepath.Join(bundle.DirPath(), pkgName)

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
			dirPath := filepath.Join(bundle.DirPath(), pkgName)

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

func init() {
	rootCmd.AddCommand(installPkgCmd)
	rootCmd.AddCommand(removePkgCmd)
	rootCmd.AddCommand(updatePkgCmd)
	rootCmd.AddCommand(listPkgCmd)

	updatePkgCmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "all packages")
}
