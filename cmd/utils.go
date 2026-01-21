package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"pea/internal/app"
	"pea/platform"

	"github.com/spf13/cobra"
)

func openEditor(cmd *cobra.Command, path string) error {
	ed := app.GetEditorConfig()

	if ed == "" {
		// Fallback to browser or default open
		if b := os.Getenv("BROWSER"); b != "" {
			c := exec.Command("bash", "-c", b+" \""+path+"\"")
			c.Stdin = os.Stdin
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.ErrOrStderr()
			if err := c.Run(); err != nil {
				return fmt.Errorf("launching BROWSER failed: %w", err)
			}
			return nil
		}

		if err := platform.BrowserImpl.OpenFile(path); err != nil {
			return fmt.Errorf("opening default editor failed: %w", err)
		}
		return nil
	}

	// Launch editor
	// Use bash -c to allow args like "code --wait"
	c := exec.Command("bash", "-c", ed+" \""+path+"\"")
	c.Stdin = os.Stdin
	c.Stdout = cmd.OutOrStdout()
	c.Stderr = cmd.ErrOrStderr()
	return c.Run()
}

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func isInputFromPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) == 0
}
