package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func addCompletionCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:       "completion [bash|zsh|install]",
		Short:     "output or install shell completion",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{"bash", "zsh", "install"},
		RunE: func(cmd *cobra.Command, args []string) error {
			op := args[0]
			switch op {
			case "bash":
				return root.GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return root.GenZshCompletion(cmd.OutOrStdout())
			case "install":
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("resolve home: %w", err)
				}
				base := filepath.Join(home, ".pea")
				if err := os.MkdirAll(base, 0o755); err != nil {
					return fmt.Errorf("create dir %s: %w", base, err)
				}

				shell := filepath.Base(os.Getenv("SHELL"))
				if shell == "zsh" {
					path := filepath.Join(base, "_pea")
					f, err := os.Create(path)
					if err != nil {
						return err
					}
					if err := root.GenZshCompletion(f); err != nil {
						f.Close()
						return err
					}
					f.Close()
					fmt.Printf("Installed completion script to %s\n", path)

					// Patch .zshrc
					rcPath := filepath.Join(home, ".zshrc")
					cfgLine := fmt.Sprintf("fpath=(%s $fpath); autoload -U compinit; compinit", base)
					if err := updateShellConfig(rcPath, cfgLine); err != nil {
						fmt.Printf("Warning: failed to update %s: %v\n", rcPath, err)
					} else {
						fmt.Printf("Updated %s\n", rcPath)
					}
					fmt.Printf("\nTo apply changes, run:\n  source %s\n", rcPath)

				} else if shell == "bash" {
					path := filepath.Join(base, "pea.bash")
					f, err := os.Create(path)
					if err != nil {
						return err
					}
					if err := root.GenBashCompletion(f); err != nil {
						f.Close()
						return err
					}
					f.Close()
					fmt.Printf("Installed completion script to %s\n", path)

					// Patch .bashrc
					rcPath := filepath.Join(home, ".bashrc")
					// Check .bash_profile on Mac if .bashrc doesn't exist?
					// stick to standard .bashrc for now.
					cfgLine := fmt.Sprintf("source %s", path)
					if err := updateShellConfig(rcPath, cfgLine); err != nil {
						fmt.Printf("Warning: failed to update %s: %v\n", rcPath, err)
					} else {
						fmt.Printf("Updated %s\n", rcPath)
					}
					fmt.Printf("\nTo apply changes, run:\n  source %s\n", rcPath)

				} else {
					// Fallback: install both, print generic
					bashPath := filepath.Join(base, "pea.bash")
					zshPath := filepath.Join(base, "_pea")

					f1, _ := os.Create(bashPath)
					root.GenBashCompletion(f1)
					f1.Close()

					f2, _ := os.Create(zshPath)
					root.GenZshCompletion(f2)
					f2.Close()

					fmt.Printf("Unknown shell '%s'. Installed both scripts to %s\n", shell, base)
					fmt.Printf("Add the relevant line to your config:\n")
					fmt.Printf("Bash: source %s\n", bashPath)
					fmt.Printf("Zsh:  fpath=(%s $fpath); autoload -U compinit; compinit\n", base)
				}
				return nil
			default:
				return cmd.Help()
			}
		},
	}
	root.AddCommand(cmd)
}

func updateShellConfig(path, command string) error {
	startMarker := "# PEA_COMPLETION_START"
	endMarker := "# PEA_COMPLETION_END"
	block := fmt.Sprintf("%s\n%s\n%s", startMarker, command, endMarker)

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(path, []byte(block+"\n"), 0o644)
		}
		return err
	}
	s := string(content)

	startIndex := strings.Index(s, startMarker)
	endIndex := strings.Index(s, endMarker)

	if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
		// Replace existing block
		before := s[:startIndex]
		after := s[endIndex+len(endMarker):]
		// Clean up potential double newlines if we remove a block
		newContent := strings.TrimRight(before, "\n") + "\n" + block + "\n" + strings.TrimLeft(after, "\n")
		return os.WriteFile(path, []byte(newContent), 0o644)
	}

	// Append
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	s += block + "\n"
	return os.WriteFile(path, []byte(s), 0o644)
}
