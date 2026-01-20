package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoGitMode(t *testing.T) {
	bin := buildBinary(t)

	home, err := os.MkdirTemp("", "pea-no-git-home-")
	if err != nil {
		t.Fatalf("mkdtemp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(home) })

	base := filepath.Join(home, ".pea")
	store := filepath.Join(base, "prompts")

	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatalf("mkdir base: %v", err)
	}

	// Write config with no_git = true
	cfg := filepath.Join(base, "config.toml")
	config := "no_git = true\n"
	if err := os.WriteFile(cfg, []byte(config), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Run 'pea add test_entry'
	add := exec.Command(bin, "add", "test_entry")
	add.Env = append(os.Environ(), "HOME="+home)
	add.Stdin = strings.NewReader("content\n")
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}

	// Verify .git does NOT exist
	if _, err := os.Stat(filepath.Join(store, ".git")); !os.IsNotExist(err) {
		t.Errorf("expected no .git directory, but found one")
	}

	// Verify file was written
	content, err := os.ReadFile(filepath.Join(store, "test_entry.md"))
	if err != nil {
		t.Fatalf("read entry: %v", err)
	}
	if string(content) != "content\n" {
		t.Errorf("entry content mismatch: got %q", string(content))
	}
}
