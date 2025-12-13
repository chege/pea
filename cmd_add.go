// Package main implements the p CLI.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func addAddCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "add <name> [file]",
		Short: "add a new entry by name, from editor, stdin, or a file",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ensureStore()
			if err != nil {
				return err
			}
			name := toSnake(args[0])
			path := filepath.Join(store, name+".txt")
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
						if err := browser.OpenFile(path); err != nil {
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
			data, err := io.ReadAll(src)
			if err != nil {
				return err
			}
			if err := os.WriteFile(path, data, 0o644); err != nil {
				return err
			}
			// git add + commit
			_ = exec.Command("bash", "-c", "cd '"+store+"' && git add '"+name+".txt' && git commit -m 'add "+name+"'").Run()
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name); err != nil {
				return err
			}
			return nil
		},
	}
	root.AddCommand(cmd)
}

var snakeRe = regexp.MustCompile(`[^a-z0-9_]+`)

func toSnake(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = snakeRe.ReplaceAllString(s, "")
	return s
}

func isInputFromPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) == 0
}
