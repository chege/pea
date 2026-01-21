package cmd

import (
	"pea/internal/app"
	"pea/platform"

	"github.com/spf13/cobra"
)

func addGetCommand(root *cobra.Command) {
	var rev string

	cmd := &cobra.Command{
		Use:               "get <name>",
		Short:             "retrieve a snippet",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeNames,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.EnsureStore()
			if err != nil {
				return err
			}

			b, err := app.ReadEntry(store, args[0], rev)
			if err != nil {
				return err
			}

			// Write to stdout
			_, err = cmd.OutOrStdout().Write(b)
			if err != nil {
				return err
			}

			// If stdout is a TTY, copy to clipboard
			if isTTY() {
				_ = copyToClipboard(string(b))
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&rev, "rev", "", "read entry content from a specific git ref")
	root.AddCommand(cmd)
}

func copyToClipboard(s string) error {
	// Use platform clipboard abstraction
	if err := platform.ClipboardImpl.Init(); err != nil {
		return err
	}

	return platform.ClipboardImpl.WriteText(s)
}

// Helper to access internal read functionality if needed by other commands (like root)
func ReadEntry(store, name, rev string) ([]byte, error) {
	return app.ReadEntry(store, name, rev)
}
