package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEditModifiesFileAndCommits(t *testing.T) {
	// 1. Setup
	bin := buildBinary(t)
	// Override PEA_STORE and HOME to a temp dir so we don't mess with user's files
	tmpHome := t.TempDir()
	storePath := filepath.Join(tmpHome, ".pea", "prompts")
	env := append(os.Environ(), "HOME="+tmpHome)

	// 2. Add an entry
	setupCmd := exec.Command(bin, "add", "my_entry")
	setupCmd.Env = env
	setupCmd.Stdin = strings.NewReader("original content")
	if out, err := setupCmd.CombinedOutput(); err != nil {
		t.Fatalf("setup add failed: %v\n%s", err, out)
	}

	// 3. Create a fake editor that appends text
	editorPath := filepath.Join(tmpHome, "fake_editor.sh")
	script := `#!/bin/bash
set -eu
file="$1"
echo " appended" >> "$file"
`
	if err := os.WriteFile(editorPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	// 4. Run 'pea edit my_entry'
	editCmd := exec.Command(bin, "edit", "my_entry")
	editCmd.Env = append(env, "EDITOR="+editorPath)
	if out, err := editCmd.CombinedOutput(); err != nil {
		t.Fatalf("edit failed: %v\n%s", err, out)
	}

	// 5. Verify content changed
	entryPath := filepath.Join(storePath, "my_entry.md")
	content, err := os.ReadFile(entryPath)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "appended") {
		t.Fatalf("content mismatch. got: %q, expected to contain 'appended'", string(content))
	}

	// 6. Verify a new commit was made
	gitCmd := exec.Command("git", "log", "--oneline")
	gitCmd.Dir = storePath
	out, err := gitCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 commits, got %d. Log:\n%s", len(lines), string(out))
	}
}

func TestEditNoChangeNoCommit(t *testing.T) {
	// 1. Setup
	bin := buildBinary(t)
	tmpHome := t.TempDir()
	storePath := filepath.Join(tmpHome, ".pea", "prompts")
	env := append(os.Environ(), "HOME="+tmpHome)

	// 2. Add an entry
	setupCmd := exec.Command(bin, "add", "fixed_entry")
	setupCmd.Env = env
	setupCmd.Stdin = strings.NewReader("static content")
	if out, err := setupCmd.CombinedOutput(); err != nil {
		t.Fatalf("setup add failed: %v\n%s", err, out)
	}

	// 3. Get commit count
	getCommitCount := func() string {
		c := exec.Command("git", "rev-list", "--count", "HEAD")
		c.Dir = storePath
		out, err := c.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(out))
	}

	initialCount := getCommitCount()

	// 4. Fake editor that does nothing
	editorPath := filepath.Join(tmpHome, "noop_editor.sh")
	script := `#!/bin/bash
exit 0
`
	if err := os.WriteFile(editorPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	// 5. Run 'pea edit fixed_entry'
	editCmd := exec.Command(bin, "edit", "fixed_entry")
	editCmd.Env = append(env, "EDITOR="+editorPath)
	if out, err := editCmd.CombinedOutput(); err != nil {
		t.Fatalf("edit failed: %v\n%s", err, out)
	}

	// 6. Verify commit count unchanged
	finalCount := getCommitCount()
	if initialCount != finalCount {
		t.Fatalf("commit count changed from %s to %s, expected no change", initialCount, finalCount)
	}
}
