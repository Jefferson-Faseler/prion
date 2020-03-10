package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
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

			_, err = git.PlainClone(dirPath, false, &git.CloneOptions{
				URL:      pkgURL,
				Depth:    1,
				Progress: os.Stdout,
			})
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
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		for _, pkgName := range args {
			dirPath := bundleDir() + "/" + pkgName
			fmt.Println("Updating " + pkgName)
			repo, err := git.PlainOpen(dirPath)
			handleError(err)

			worktree, err := repo.Worktree()
			handleError(err)

			err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
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

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "a subcommand for configuring vimrc",
}

// configAddCmd represents the command for adding vim configuration
var configAddCmd = &cobra.Command{
	Use:   "add [configuration]",
	Short: "easily add a vim configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		configurationText := args[0]
		fmt.Println("Adding " + configurationText + " to vimrc")

		vimrc, err := os.OpenFile(vimrcPath(), os.O_APPEND|os.O_WRONLY, 0666)
		handleError(err)

		defer vimrc.Close()
		_, err = vimrc.WriteString(configurationText)
		handleError(err)
	},
}

// configEditCmd represents the command for editing vim configuration
var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit your vim configuration manually",
	RunE: func(cmd *cobra.Command, args []string) error {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // I mean, this makes sense, right?
		}

		fmt.Println("Opening editor for manual changes. Using editor " + editor)
		fmt.Println("Close editor to continue")
		executable, err := exec.LookPath(editor)
		handleError(err)

		editorCmd := exec.Command(executable, vimrcPath())
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		return editorCmd.Run()
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

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func bundleDir() string {
	return viper.GetString("VIM_BUNDLE_DIR")
}

func vimrcPath() string {
	return viper.GetString("VIMRC_PATH")
}

func init() {
	rootCmd.AddCommand(installPkgCmd)
	rootCmd.AddCommand(removePkgCmd)
	rootCmd.AddCommand(updatePkgCmd)
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configEditCmd)
}
