package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func addAddCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "add <name> [file]",
		Short: "add a new entry by name, from editor, stdin, or a file",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ensureStore()
			if err != nil { return err }
			name := toSnake(args[0])
			path := filepath.Join(store, name+".txt")
			var src io.Reader
			if len(args) > 1 {
				f, err := os.Open(args[1])
				if err != nil { return err }
				defer f.Close()
				src = f
			} else {
				// If stdin has data, read it; else require EDITOR and open it is non-trivial, so for v0 read stdin or error.
				if isInputFromPipe() {
					src = bufio.NewReader(os.Stdin)
				} else {
					return fmt.Errorf("no input provided: supply a file or pipe stdin")
				}
			}
			data, err := io.ReadAll(src)
			if err != nil { return err }
			if err := os.WriteFile(path, data, 0o644); err != nil { return err }
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name)
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
	if err != nil { return false }
	return (fi.Mode() & os.ModeCharDevice) == 0
}
