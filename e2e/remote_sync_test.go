package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoteSyncPushesWhenConfigured(t *testing.T) {
	bin := buildBinary(t)

	home, err := os.MkdirTemp("", "pea-remote-home-")
	if err != nil {
		t.Fatalf("mkdtemp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(home) })

	base := filepath.Join(home, ".pea")
	store := filepath.Join(base, "prompts")
	remote := filepath.Join(home, "remote.git")

	if out, err := exec.Command("git", "init", "--bare", remote).CombinedOutput(); err != nil {
		t.Fatalf("init bare: %v\n%s", err, out)
	}

	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatalf("mkdir base: %v", err)
	}

	cfg := filepath.Join(base, "config.toml")
	config := "store_dir = \"" + store + "\"\nremote_url = \"" + remote + "\"\n"
	if err := os.WriteFile(cfg, []byte(config), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	add := exec.Command(bin, "add", "remote_entry")
	add.Env = append(os.Environ(), "HOME="+home)
	add.Stdin = strings.NewReader("content\n")
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}

	// Verify remote has the commit
	revParse := exec.Command("git", "--git-dir", remote, "rev-parse", "HEAD")
	if out, err := revParse.CombinedOutput(); err != nil {
		t.Fatalf("remote missing commit: %v\n%s", err, out)
	}
}

func TestRemoteCommand(t *testing.T) {
	bin := buildBinary(t)

	home, err := os.MkdirTemp("", "pea-remote-cmd-home-")
	if err != nil {
		t.Fatalf("mkdtemp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(home) })

	base := filepath.Join(home, ".pea")
	store := filepath.Join(base, "prompts")
	remote := filepath.Join(home, "remote.git")

	// Init bare remote
	if out, err := exec.Command("git", "init", "--bare", remote).CombinedOutput(); err != nil {
		t.Fatalf("init bare: %v\n%s", err, out)
	}

	// Run 'pea remote <url>'
	cmd := exec.Command(bin, "remote", remote)
	cmd.Env = append(os.Environ(), "HOME="+home)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("pea remote failed: %v\n%s", err, out)
	}

	// Verify git remote in store
	remoteCmd := exec.Command("git", "remote", "get-url", "origin")
	remoteCmd.Dir = store
	out, err := remoteCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git remote get-url failed: %v", err)
	}
	if strings.TrimSpace(string(out)) != remote {
		t.Errorf("git remote url mismatch: got %q, want %q", string(out), remote)
	}
}
