package cmd

import (
	"fmt"
	"strings"

	"github.com/Jefferson-Faseler/prion/internal/bundle"
	"github.com/Jefferson-Faseler/prion/pkgmngr"

	"github.com/spf13/cobra"
)

var (
	all bool // used for updating all packages
)

// installPkgCmd represents the command for installing vim packages
var installPkgCmd = &cobra.Command{
	Use:           "install [pkg url]",
	Aliases:       []string{"add", "i"},
	Short:         "easily install a vim package",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return nil
		}

		for _, pkgURL := range args {
			fmt.Println("Installing " + pkgURL)
			err := pkgmngr.Install(pkgURL)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

// ReinstallPkgCmd represents the command for installing vim packages
var reinstallPkgCmd = &cobra.Command{
	Use:           "reinstall [pkg url]",
	Short:         "reinstall a vim package",
	Long:          "Use when the package cannot update with a simple pull",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return nil
		}

		for _, pkgName := range args {
			fmt.Println("Reinstalling " + pkgName)
			err := pkgmngr.Reinstall(pkgName)
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
	Aliases:       []string{"remove", "delete"},
	Short:         "easily remove a vim package",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return nil
		}

		for _, pkgName := range args {
			err := pkgmngr.Remove(pkgName)
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
			pkgs, err = bundle.Packages()
			if err != nil {
				return err
			}
		} else {
			pkgs = args
		}

		for _, pkgName := range pkgs {
			fmt.Println("Updating " + pkgName)
			wasUpToDate, err := pkgmngr.Update(pkgName)
			if err != nil {
				if strings.Contains(err.Error(), "object not found") ||
					strings.Contains(err.Error(), "ssh: handshake failed: knownhosts: key mismatch") {
					fmt.Println(err.Error())
					continue
				}
				return err
			}
			if wasUpToDate {
				fmt.Println("already up-to-date")
				continue
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
	rootCmd.AddCommand(reinstallPkgCmd)
	rootCmd.AddCommand(removePkgCmd)
	rootCmd.AddCommand(updatePkgCmd)
	rootCmd.AddCommand(listPkgCmd)

	updatePkgCmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "all packages")
}
