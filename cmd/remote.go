package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"pea/internal/app"
	"strings"

	"github.com/spf13/cobra"
)

func addRemoteCommand(root *cobra.Command) {
	remoteCmd := &cobra.Command{
		Use:   "remote <url> | create <name>",
		Short: "Configure a remote git repository",
		Long:  "Configure a remote git repository to sync your prompts.\nProvide a URL to set an existing remote, or use 'create' to make a new one via GitHub CLI.",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}
			// Existing behavior: pea remote <url>
			return runSetRemote(cmd, args[0])
		},
	}

	createCmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new GitHub repository via 'gh' CLI",
		Long:  "Create a new repository on GitHub using the 'gh' CLI and configure it as the remote.\nRequires 'gh' to be installed and authenticated.",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreateRemote,
	}
	createCmd.Flags().Bool("public", false, "make the new repository public (default: private)")

	remoteCmd.AddCommand(createCmd)
	root.AddCommand(remoteCmd)
}

func runSetRemote(cmd *cobra.Command, url string) error {
	// 1. Ensure store is loaded and get path
	store, err := app.EnsureStore()
	if err != nil {
		return err
	}

	// 2. Configure git remote (force update)
	if err := app.SetGitRemote(store, url); err != nil {
		return err
	}

	fmt.Printf("Remote configured: %s\n", url)
	fmt.Println("Future changes will be automatically pushed.")
	return nil
}

func runCreateRemote(cmd *cobra.Command, args []string) error {
	repoName := args[0]
	isPublic, _ := cmd.Flags().GetBool("public")

	// 1. Check gh installed
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("GitHub CLI 'gh' not found. Please install it: https://cli.github.com/")
	}

	// 2. Check login status
	out, err := exec.Command("gh", "api", "user", "--jq", ".login").Output()
	if err != nil {
		return fmt.Errorf("failed to check GitHub login status. Run 'gh auth login' first.\nError: %w", err)
	}
	user := strings.TrimSpace(string(out))
	if user == "" {
		return fmt.Errorf("not logged in to GitHub. Run 'gh auth login' first.")
	}

	// 3. Confirm
	visibility := "private"
	if isPublic {
		visibility = "public"
	}
	fmt.Printf("Logged in as: %s\n", user)
	fmt.Printf("Create and sync with 'github.com/%s/%s' (%s)? [y/N] ", user, repoName, visibility)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		fmt.Println("Aborted.")
		return nil
	}

	// 4. Ensure store exists
	store, err := app.EnsureStore()
	if err != nil {
		return err
	}

	// 5. Run gh repo create
	// gh repo create <name> --private --source=. --remote=origin --push
	ghArgs := []string{
		"repo", "create", repoName,
		"--source=.",
		"--remote=origin",
		"--push",
	}
	if isPublic {
		ghArgs = append(ghArgs, "--public")
	} else {
		ghArgs = append(ghArgs, "--private")
	}

	fmt.Printf("Creating repository...\n")
	ghCmd := exec.Command("gh", ghArgs...)
	ghCmd.Dir = store
	ghCmd.Stdout = cmd.OutOrStdout()
	ghCmd.Stderr = cmd.ErrOrStderr()

	if err := ghCmd.Run(); err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// 6. Get the new remote URL (optional check)
	remoteURLBytes, err := exec.Command("git", "-C", store, "remote", "get-url", "origin").Output()
	if err != nil {
		fmt.Printf("Warning: could not retrieve remote URL from git: %v\n", err)
	} else {
		// Just print it for verification
		fmt.Printf("Linked to remote: %s\n", strings.TrimSpace(string(remoteURLBytes)))
	}

	fmt.Println("\n✅ Repository created and linked successfully!")
	return nil
}
