package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// gitAddAndCommit attempts to stage the provided paths and commit with commitMsg.
// It logs failures to stderr but does not block the calling command.
func gitAddAndCommit(store string, paths []string, commitMsg string, stderr io.Writer) {
	if stderr == nil {
		stderr = io.Discard
	}

	addArgs := append([]string{"add"}, paths...)
	add := exec.Command("git", addArgs...)
	add.Dir = store
	if out, err := add.CombinedOutput(); err != nil {
		fmt.Fprintf(stderr, "warning: git add failed: %v: %s\n", err, string(out))
		return
	}

	commit := exec.Command("git", "commit", "-m", commitMsg)
	commit.Dir = store
	if out, err := commit.CombinedOutput(); err != nil {
		fmt.Fprintf(stderr, "warning: git commit failed: %v: %s\n", err, string(out))
	}
}

// gitRmAndCommit attempts to stage deletions and commit with commitMsg.
// It logs failures to stderr but does not block the calling command.
func gitRmAndCommit(store string, paths []string, commitMsg string, stderr io.Writer) {
	if stderr == nil {
		stderr = io.Discard
	}

	rmArgs := append([]string{"rm", "-f"}, paths...)
	rm := exec.Command("git", rmArgs...)
	rm.Dir = store
	if out, err := rm.CombinedOutput(); err != nil {
		fmt.Fprintf(stderr, "warning: git rm failed: %v: %s\n", err, string(out))
		return
	}

	commit := exec.Command("git", "commit", "-m", commitMsg)
	commit.Dir = store
	if out, err := commit.CombinedOutput(); err != nil {
		fmt.Fprintf(stderr, "warning: git commit failed: %v: %s\n", err, string(out))
	}
}

func revertLastCommitForPath(store, path string, stderr io.Writer) error {
	if stderr == nil {
		stderr = io.Discard
	}

	logCmd := exec.Command("git", "log", "-n1", "--format=%H", "--", path)
	logCmd.Dir = store
	shaOut, err := logCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git log failed: %w: %s", err, string(shaOut))
	}
	sha := strings.TrimSpace(string(shaOut))
	if sha == "" {
		return fmt.Errorf("no commits found for %s", path)
	}
	revert := exec.Command("git", "revert", "--no-edit", sha)
	revert.Dir = store
	if out, err := revert.CombinedOutput(); err != nil {
		return fmt.Errorf("git revert failed: %w: %s", err, string(out))
	}
	return nil
}
