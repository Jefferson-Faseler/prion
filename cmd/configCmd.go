package cmd

import (
	"fmt"

	"github.com/Jefferson-Faseler/prion/internal/vimconfig"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "a subcommand for configuring vimrc",
}

// configAddCmd represents the command for adding vim configuration
var configAddCmd = &cobra.Command{
	Use:           "add [configuration]",
	Short:         "easily add a vim configuration",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return nil
		}

		configurationText := args[0]
		fmt.Println("Adding " + configurationText + " to vimrc")

		vimrc, err := vimconfig.PrepareForWrite()
		if err != nil {
			return err
		}

		defer vimrc.Close()
		return vimrc.Write(configurationText)
	},
}

// configEditCmd represents the command for editing vim configuration
var configEditCmd = &cobra.Command{
	Use:           "edit",
	Short:         "edit your vim configuration manually",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Opening editor for changes.")
		return vimconfig.EditVimRC()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configEditCmd)
}
