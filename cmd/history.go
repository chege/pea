package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"pea/internal/app"
	"strconv"

	"github.com/spf13/cobra"
)

func addHistoryCommand(root *cobra.Command) {
	var limit int
	var reverse bool

	cmd := &cobra.Command{
		Use:   "history <name>",
		Short: "show git history for an entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.EnsureStore()
			if err != nil {
				return err
			}

			name, err := app.NormalizeName(args[0])
			if err != nil {
				return err
			}

			_, ext, err := app.ExistingEntryPath(store, name)
			if err != nil {
				return fmt.Errorf("history failed: not found: %s", name)
			}

			path := name + ext

			argsLog := []string{"log", "--follow", "--pretty=format:%h %s", "--max-count", strconv.Itoa(limit)}
			if reverse {
				argsLog = append(argsLog, "--reverse")
			}
			argsLog = append(argsLog, "--", path)

			c := exec.Command("git", argsLog...)
			c.Dir = store
			var out bytes.Buffer
			c.Stdout = &out
			c.Stderr = &out
			if err := c.Run(); err != nil {
				return fmt.Errorf("history failed: %w: %s", err, out.String())
			}

			if out.Len() == 0 {
				_, err := fmt.Fprintf(cmd.OutOrStdout(), "no history for %s\n", name)
				return err
			}

			if _, err := cmd.OutOrStdout().Write(out.Bytes()); err != nil {
				return err
			}
			// Ensure trailing newline
			if out.Bytes()[out.Len()-1] != '\n' {
				_, _ = fmt.Fprintln(cmd.OutOrStdout())
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 20, "maximum number of entries to show")
	cmd.Flags().BoolVar(&reverse, "reverse", false, "show oldest first")
	root.AddCommand(cmd)
}
