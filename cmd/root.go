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
	base, defaultStore := app.DefaultPaths()
	cfgPath := filepath.Join(base, "config.toml")
	cmd := &cobra.Command{
		Use:   "pea",
		Short: "pea: fast local prompt storage & retrieval",
		Long:  fmt.Sprintf("pea is a fast, local CLI to store short text under names and retrieve it instantly.\n\nDefaults: store at %s; config at %s; env override: PEA_STORE (highest precedence).", defaultStore, cfgPath),
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = false

	cmd.Version = version
	cmd.SetVersionTemplate("pea version {{.Version}}\n")

	addListCommand(cmd)
	addGetCommand(cmd)
	addCpCommand(cmd)
	addAddCommand(cmd)
	addEditCommand(cmd)
	addRemoveCommand(cmd)
	addMoveCommand(cmd)
	addHistoryCommand(cmd)
	addSearchCommand(cmd)
	addCompletionCommand(cmd)
	addRemoteCommand(cmd)
	addSyncCommand(cmd)
	return cmd
}

func Execute() {
	root := NewRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
