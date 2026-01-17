package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

func newRootCmd() *cobra.Command {
	var rev string
	base, defaultStore := defaultPaths()
	cfgPath := filepath.Join(base, "config.toml")
	cmd := &cobra.Command{
		Use:   "pea [name]",
		Short: "pea: fast local prompt storage & retrieval",
		Long:  fmt.Sprintf("pea is a fast, local CLI to store short text under names and retrieve it instantly.\n\nDefaults: store at %s; config at %s; env override: PEA_STORE (highest precedence).", defaultStore, cfgPath),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 1 {

				store, err := ensureStore()

				if err != nil {
					return err
				}

				b, err := readEntry(store, args[0], rev)

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

func main() {
	root := newRootCmd()

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
