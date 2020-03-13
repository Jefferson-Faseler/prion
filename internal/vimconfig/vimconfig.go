package vimconfig

import (
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

// VimRC is an alias of os.File
type VimRC struct {
	*os.File
}

// PrepareForWrite opens a file for writing
func PrepareForWrite() (*VimRC, error) {
	f, err := os.OpenFile(vimrcPath(), os.O_APPEND|os.O_WRONLY, 0666)
	return &VimRC{f}, err
}

// Write writes the provided cnfiguration text to a file
func (v *VimRC) Write(configurationText string) error {
	_, err := v.File.WriteString(configurationText)
	return err
}

// EditVimRC opens the default editor defined by the EDITOR environment variable for a user. If none is set it defaults to vim.
func EditVimRC() error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // I mean, this makes sense, right?
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	editorCmd := exec.Command(executable, vimrcPath())
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	return editorCmd.Run()
}

func vimrcPath() string {
	return viper.GetString("VIMRC_PATH")
}
