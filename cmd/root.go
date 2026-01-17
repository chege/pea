package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"pea/internal/app"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

func NewRootCmd() *cobra.Command {
	var rev string
	base, defaultStore := app.DefaultPaths()
	cfgPath := filepath.Join(base, "config.toml")
	cmd := &cobra.Command{
		Use:               "pea [name]",
		Short:             "pea: fast local prompt storage & retrieval",
		Long:              fmt.Sprintf("pea is a fast, local CLI to store short text under names and retrieve it instantly.\n\nDefaults: store at %s; config at %s; env override: PEA_STORE (highest precedence).", defaultStore, cfgPath),
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: completeNames,
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 1 {

				store, err := app.EnsureStore()

				if err != nil {
					return err
				}

				// ReadEntry is defined in get.go (helper wrapping app.ReadEntry)
				// or we can call app.ReadEntry directly.
				// Since we need isTTY and copyToClipboard which are in get.go (package cmd),
				// we are good.
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
			}

			return cmd.Help()
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = false

	cmd.Version = version
	cmd.SetVersionTemplate("pea version {{.Version}}\n")

	addListCommand(cmd)
	addAddCommand(cmd)
	addRemoveCommand(cmd)
	addMoveCommand(cmd)
	addHistoryCommand(cmd)
	addSearchCommand(cmd)
	addCompletionCommand(cmd)
	cmd.Flags().StringVar(&rev, "rev", "", "read entry content from a specific git ref")
	return cmd
}

func Execute() {
	root := NewRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
