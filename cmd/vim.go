package cmd

import (
	"os"

	"gopkg.in/src-d/go-git.v4"

	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// vimCmd represents the vim subcommand
var vimCmd = &cobra.Command{
	Use:   "vim",
	Short: "a subcommand for configuring vim",
}

// vimAddPkgCmd represents the command for adding vim packages
var vimAddPkgCmd = &cobra.Command{
	Use:   "addpkg [pkg url]",
	Short: "easily add vim packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgURL := args[0]

		log.Print("Adding " + pkgURL + " to vim packages")

		//location := viper.GetString("VIM_PACKAGE_DIR") + "/"
		//log.Print(location)
		_, err := git.PlainClone("", false, &git.CloneOptions{
			URL:      pkgURL,
			Depth:    1,
			Progress: os.Stdout,
		})

		if err != nil {
			log.Fatal(err)
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(vimCmd)
	vimCmd.AddCommand(vimAddPkgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vimCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vimCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
