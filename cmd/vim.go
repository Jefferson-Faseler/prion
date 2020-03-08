package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
)

// vimCmd represents the vim subcommand
var vimCmd = &cobra.Command{
	Use:   "vim",
	Short: "a subcommand for configuring vim",
}

var vimPkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "a subcommand for working with packages",
}

// vimAddPkgCmd represents the command for adding vim packages
var vimAddPkgCmd = &cobra.Command{
	Use:   "add [pkg url]",
	Short: "easily add vim packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		pkgURL := args[0]

		pkgName := getPkgName(pkgURL)
		dirPath := viper.GetString("VIM_PACKAGE_DIR") + "/" + pkgName
		log.Print("Adding " + pkgURL + " to vim packages at " + dirPath)
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}

		_, err = git.PlainClone(dirPath, false, &git.CloneOptions{
			URL:      pkgURL,
			Depth:    1,
			Progress: os.Stdout,
		})
		if err != nil {
			log.Panic(err)
		}
		return err
	},
}

// vimRemovePkgCmd represents the command for adding vim packages
var vimRemovePkgCmd = &cobra.Command{
	Use:   "rm [pkg name]",
	Short: "easily remove vim packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		pkgName := args[0]

		dirPath := viper.GetString("VIM_PACKAGE_DIR") + "/" + pkgName
		log.Print("Removing " + pkgName + " from vim packages found in " + dirPath)
		err := os.RemoveAll(dirPath)
		if err != nil {
			log.Panic(err)
		}
		return err
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

func init() {
	rootCmd.AddCommand(vimCmd)
	vimCmd.AddCommand(vimPkgCmd)
	vimPkgCmd.AddCommand(vimAddPkgCmd)
	vimPkgCmd.AddCommand(vimRemovePkgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vimCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vimCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
