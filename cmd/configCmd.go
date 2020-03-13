package cmd

import (
	"fmt"
	"os"

	"github.com/Jefferson-Faseler/prion/internal/vimconfig"
	"github.com/spf13/cobra"
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

		vimrc, err := vimconfig.PrepareForWrite()
		handleError(err)

		defer vimrc.Close()
		err = vimrc.Write(configurationText)
		handleError(err)
	},
}

// configEditCmd represents the command for editing vim configuration
var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit your vim configuration manually",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Opening editor for changes.")
		err := vimconfig.EditVimRC()
		handleError(err)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configEditCmd)
}
