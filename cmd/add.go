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

	"github.com/spf13/cobra"
)

func addAddCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "add <name> [file]",
		Short: "add a new entry by name, from editor, stdin, or a file",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.EnsureStore()

			if err != nil {
				return err
			}

			name, err := app.NormalizeName(args[0])
			if err != nil {
				return err
			}

			path, ext, err := app.TargetEntryPath(store, name)
			if err != nil {
				return err
			}

			var src io.Reader

			if len(args) > 1 {
				f, err := os.Open(args[1])

				if err != nil {
					return err
				}

				defer func() { _ = f.Close() }()
				src = f
			} else {
				// If stdin has data, read it; else open $EDITOR for the target path
				if isInputFromPipe() {
					src = bufio.NewReader(os.Stdin)
				} else {

					ed := os.Getenv("EDITOR")

					// Ensure file exists before opening editor
					if _, err := os.Stat(path); os.IsNotExist(err) {
						if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
							return err
						}
					} else if err != nil {
						return err
					}

					if ed == "" {
						if b := os.Getenv("BROWSER"); b != "" {
							c := exec.Command("bash", "-c", b+" \""+path+"\"")
							c.Stdin = os.Stdin
							c.Stdout = cmd.OutOrStdout()
							c.Stderr = cmd.ErrOrStderr()
							if err := c.Run(); err != nil {
								return fmt.Errorf("$EDITOR is not set and launching BROWSER failed: %w", err)
							}
							if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name); err != nil {
								return err
							}
							return nil
						}

						if err := platform.BrowserImpl.OpenFile(path); err != nil {
							return fmt.Errorf("$EDITOR is not set and opening default editor failed: %w", err)
						}

						if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name); err != nil {
							return err
						}

						return nil
					}

					// Launch editor; content handled by editor, then print name and exit
					c := exec.Command("bash", "-c", ed+" \""+path+"\"")
					c.Stdin = os.Stdin
					c.Stdout = cmd.OutOrStdout()
					c.Stderr = cmd.ErrOrStderr()
					if err := c.Run(); err != nil {
						return err
					}

					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name); err != nil {
						return err
					}

					return nil
				}
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

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name); err != nil {
				return err
			}

			return nil
		},
	}
	root.AddCommand(cmd)
}

func isInputFromPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) == 0
}
