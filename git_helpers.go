package main

import (
	"fmt"
	"io"
	"os/exec"
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
