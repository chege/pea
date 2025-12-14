package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnvOverridesConfig(t *testing.T) {
	bin := buildBinary(t)

	home := os.Getenv("HOME")
	base := filepath.Join(home, ".pea")
	configStore := filepath.Join(base, "config_store")
	envStore := filepath.Join(base, "env_store")

	_ = os.MkdirAll(configStore, 0o755)
	_ = os.MkdirAll(envStore, 0o755)

	cfg := filepath.Join(base, "config.toml")
	if err := os.WriteFile(cfg, []byte("store_dir = \""+configStore+"\"\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(cfg)
		_ = os.RemoveAll(configStore)
		_ = os.RemoveAll(envStore)
	})

	cmd := exec.Command(bin, "add", "env_entry")
	cmd.Env = append(os.Environ(), "PEA_STORE="+envStore)
	cmd.Stdin = strings.NewReader("env store\n")

	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}

	if _, err := os.Stat(filepath.Join(envStore, "env_entry.md")); err != nil {
		t.Fatalf("entry not in env store: %v", err)
	}

	if _, err := os.Stat(filepath.Join(configStore, "env_entry.md")); err == nil {
		t.Fatalf("entry should not be in config store")
	}
}

func TestConfigRejectsRelativeStoreDir(t *testing.T) {
	bin := buildBinary(t)

	home := os.Getenv("HOME")
	base := filepath.Join(home, ".pea")
	_ = os.MkdirAll(base, 0o755)
	cfg := filepath.Join(base, "config.toml")

	if err := os.WriteFile(cfg, []byte("store_dir = \"relative/path\"\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(cfg)
	})

	cmd := exec.Command(bin, "ls")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected ls to fail for relative store_dir, output: %s", out)
	}

	if !strings.Contains(string(out), "store_dir must be an absolute path") {
		t.Fatalf("expected error about absolute path, got: %s", out)
	}
}
