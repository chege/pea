package cmd

import (
	"fmt"
	"pea/internal/app"
	"pea/platform"

	"github.com/spf13/cobra"
)

func addCpCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:               "cp <name>",
		Short:             "copy a snippet to the clipboard",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeNames,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.EnsureStore()
			if err != nil {
				return err
			}

			// Read content
			b, err := app.ReadEntry(store, args[0], "")
			if err != nil {
				return err
			}

			// Copy to clipboard
			if err := platform.ClipboardImpl.Init(); err != nil {
				return fmt.Errorf("clipboard init failed: %w", err)
			}
			if err := platform.ClipboardImpl.WriteText(string(b)); err != nil {
				return fmt.Errorf("clipboard write failed: %w", err)
			}

			fmt.Fprintf(cmd.ErrOrStderr(), "✓ Copied '%s' to clipboard.\n", args[0])
			return nil
		},
	}
	root.AddCommand(cmd)
}
