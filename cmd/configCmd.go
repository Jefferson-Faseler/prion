package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

func vimrcPath() string {
	return viper.GetString("VIMRC_PATH")
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configEditCmd)
}
