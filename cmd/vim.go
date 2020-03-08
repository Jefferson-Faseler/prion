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
	Short: "easily add a vim package",
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
	Short: "easily remove a vim package",
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

// vimUpdatePkgCmd represents the command for adding vim packages
var vimUpdatePkgCmd = &cobra.Command{
	Use:   "update [pkg url]",
	Short: "easily update a vim package",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		pkgName := args[0]

		dirPath := viper.GetString("VIM_PACKAGE_DIR") + "/" + pkgName
		log.Print("Updating " + pkgName + " in vim packages found at " + dirPath)
		repo, err := git.PlainOpen(dirPath)
		if err != nil {
			log.Panic(err)
		}
		worktree, err := repo.Worktree()
		if err != nil {
			log.Panic(err)
		}
		err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			if strings.Contains(err.Error(), "already up-to-date") {
				log.Print(err)
				os.Exit(0)
			}
			log.Panic(err)
		}
		ref, err := repo.Head()
		if err != nil {
			log.Panic(err)
		}
		log.Print(ref)
		return err
	},
}

var vimConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "a subcommand for configuring vimrc",
}

// vimAddConfigCmd represents the command for adding vim packages
var vimAddConfigCmd = &cobra.Command{
	Use:   "add [configuration]",
	Short: "easily add a vim configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		configurationText := args[0]
		vimrcPath := viper.GetString("VIMRC_PATH")
		log.Print("Adding " + configurationText + " to vimrc " + vimrcPath)

		vimrc, err := os.OpenFile(vimrcPath, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Panic(err)
		}
		defer vimrc.Close()
		if _, err = vimrc.WriteString(configurationText); err != nil {
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
	vimCmd.AddCommand(vimConfigCmd)
	vimPkgCmd.AddCommand(vimAddPkgCmd)
	vimPkgCmd.AddCommand(vimRemovePkgCmd)
	vimPkgCmd.AddCommand(vimUpdatePkgCmd)
	vimConfigCmd.AddCommand(vimAddConfigCmd)
}
