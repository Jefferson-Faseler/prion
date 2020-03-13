package cmd

import (
	"fmt"
	"os"
	"strings"

	"prion/internal/bundle"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	all bool // used for updating all packages
)

// installPkgCmd represents the command for installing vim packages
var installPkgCmd = &cobra.Command{
	Use:   "install [pkg url]",
	Short: "easily install a vim package",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		for _, pkgURL := range args {
			pkgName := getPkgName(pkgURL)
			dirPath := bundleDir() + "/" + pkgName
			fmt.Println("Installing " + pkgURL)
			err := os.MkdirAll(dirPath, os.ModePerm)
			handleError(err)

			err = bundle.Clone(pkgURL, dirPath)
			if err != nil {
				if nestedErr := os.RemoveAll(dirPath); nestedErr != nil {
					fmt.Println(err.Error())
				}
				handleError(err)
			}
		}
	},
}

// removePkgCmd represents the command for removing vim packages
var removePkgCmd = &cobra.Command{
	Use:   "rm [pkg name]",
	Short: "easily remove a vim package",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		for _, pkgName := range args {
			dirPath := bundleDir() + "/" + pkgName
			err := os.RemoveAll(dirPath)
			handleError(err)
			fmt.Println(pkgName + " removed")
		}
	},
}

// updatePkgCmd represents the command for updating vim packages
var updatePkgCmd = &cobra.Command{
	Use:   "update [pkg url]",
	Short: "easily update a vim package",
	Run: func(cmd *cobra.Command, args []string) {
		var pkgs []string
		var err error

		if len(args) == 0 && all == false {
			cmd.Help()
			os.Exit(0)
		} else if all == true {
			_pkgs, err := bundle.Packages()
			pkgs = _pkgs
			handleError(err)
		} else {
			pkgs = args
		}

		for _, pkgName := range pkgs {
			dirPath := bundleDir() + "/" + pkgName
			fmt.Println("Updating " + pkgName)

			err = bundle.Pull(dirPath)
			if err != nil {
				if strings.Contains(err.Error(), "already up-to-date") {
					fmt.Println(err)
				} else {
					handleError(err)
				}
			}
		}
	},
}

// listPkgCmd represents the command for updating vim packages
var listPkgCmd = &cobra.Command{
	Use:   "ls [pkg url]",
	Short: "list all your vim packages",
	Run: func(cmd *cobra.Command, args []string) {
		pkgs, err := bundle.Packages()
		handleError(err)
		for _, pkgName := range pkgs {
			fmt.Println(pkgName)
		}
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

func init() {
	rootCmd.AddCommand(installPkgCmd)
	rootCmd.AddCommand(removePkgCmd)
	rootCmd.AddCommand(updatePkgCmd)
	rootCmd.AddCommand(listPkgCmd)

	updatePkgCmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "all packages")
}
