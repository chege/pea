package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGitDisabledByConfig(t *testing.T) {
	bin := buildBinary(t)

	home, err := os.MkdirTemp("", "pea-no-git-conf-")
	if err != nil {
		t.Fatalf("mkdtemp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(home) })

	base := filepath.Join(home, ".pea")
	store := filepath.Join(base, "prompts")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatalf("mkdir base: %v", err)
	}

	// Explicitly disable git in config
	cfg := filepath.Join(base, "config.toml")
	config := "git = false\n"
	if err := os.WriteFile(cfg, []byte(config), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Run add
	add := exec.Command(bin, "add", "no_git_entry")
	add.Env = append(os.Environ(), "HOME="+home)
	add.Stdin = strings.NewReader("content\n")
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}

	// Verify .git does NOT exist
	if _, err := os.Stat(filepath.Join(store, ".git")); !os.IsNotExist(err) {
		t.Errorf("expected no .git directory when git=false, but found one")
	}
}

func TestGitAutoSkipIfMissing(t *testing.T) {
	bin := buildBinary(t)

	home, err := os.MkdirTemp("", "pea-missing-git-")
	if err != nil {
		t.Fatalf("mkdtemp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(home) })

	base := filepath.Join(home, ".pea")
	store := filepath.Join(base, "prompts")

	// Run add with PATH cleared so 'git' is not found
	add := exec.Command(bin, "add", "missing_git_entry")
	add.Env = []string{"HOME=" + home, "PATH="} // Empty PATH
	add.Stdin = strings.NewReader("content\n")

	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("add failed (should have succeeded without git): %v\n%s", err, out)
	}

	// Verify content exists
	if _, err := os.Stat(filepath.Join(store, "missing_git_entry.md")); err != nil {
		t.Errorf("entry not created")
	}

	// Verify .git does NOT exist
	if _, err := os.Stat(filepath.Join(store, ".git")); !os.IsNotExist(err) {
		t.Errorf("expected no .git directory when git binary missing, but found one")
	}
}
