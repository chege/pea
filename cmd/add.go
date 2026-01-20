package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"pea/internal/app"
	"pea/platform"
	"strings"

	"github.com/spf13/cobra"
)

func addAddCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "add [name] [file]",
		Short: "add a new entry by name, from editor, stdin, or a file",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return runAddInteractive(cmd)
			}
			return runAddNamed(cmd, args)
		},
	}
	root.AddCommand(cmd)
}

func runAddInteractive(cmd *cobra.Command) error {
	store, err := app.EnsureStore()
	if err != nil {
		return err
	}

	// Create a temp file
	tmpFile, err := os.CreateTemp("", "pea-*.md")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Open editor
	if err := openEditor(cmd, tmpPath); err != nil {
		return err
	}

	// Read content
	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return fmt.Errorf("add aborted: empty content")
	}

	// Ask for name
	fmt.Fprintf(cmd.OutOrStdout(), "Enter name: ")
	reader := bufio.NewReader(os.Stdin)
	nameInput, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read name: %w", err)
	}
	name := strings.TrimSpace(nameInput)
	if name == "" {
		return fmt.Errorf("add aborted: empty name")
	}

	return saveEntry(cmd, store, name, bytes.NewReader(data))
}

func runAddNamed(cmd *cobra.Command, args []string) error {
	store, err := app.EnsureStore()
	if err != nil {
		return err
	}

	name, err := app.NormalizeName(args[0])
	if err != nil {
		return err
	}

	path, _, err := app.TargetEntryPath(store, name)
	if err != nil {
		return err
	}

	var src io.Reader

	if len(args) > 1 {
		f, err := os.Open(args[1])
		if err != nil {
			return err
		}
		defer f.Close()
		src = f
	} else {
		if isInputFromPipe() {
			src = bufio.NewReader(os.Stdin)
		} else {
			// Ensure file exists before opening editor
			if _, err := os.Stat(path); os.IsNotExist(err) {
				if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}

			if err := openEditor(cmd, path); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name)
			return nil
		}
	}

	return saveEntry(cmd, store, name, src)
}

func saveEntry(cmd *cobra.Command, store, name string, src io.Reader) error {
	path, ext, err := app.TargetEntryPath(store, name)
	if err != nil {
		return err
	}

	existedBefore := app.FileExists(path)
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	if len(bytes.TrimSpace(data)) == 0 {
		if !existedBefore {
			_ = os.Remove(path)
		}
		return fmt.Errorf("add failed: empty content")
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}

	// git add + commit (best-effort)
	commitMsg := "feat: add " + name + ext
	app.GitAddAndCommit(store, []string{name + ext}, commitMsg, cmd.ErrOrStderr())

	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name)
	return nil
}

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

func isInputFromPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) == 0
}
