package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestTomlStoreDir(t *testing.T) {

	bin := buildBinary(t)

	// Write config.toml with store_dir
	home, _ := os.UserHomeDir()

	base := filepath.Join(home, ".pea")

	_ = os.MkdirAll(base, 0o755)

	store := filepath.Join(base, "toml_store")

	_ = os.MkdirAll(store, 0o755)

	cfg := filepath.Join(base, "config.toml")

	if err := os.WriteFile(cfg, []byte("store_dir = \""+store+"\"\n"), 0o644); err != nil {

		t.Fatal(err)
	}

	// Add should write to toml store
	cmd := exec.Command(bin, "add", "toml_entry")

	cmd.Stdin = strings.NewReader("conf\n")

	if out, err := cmd.CombinedOutput(); err != nil {

		t.Fatalf("add failed: %v\n%s", err, out)
	}

	if _, err := os.Stat(filepath.Join(store, "toml_entry.txt")); err != nil {

		t.Fatalf("entry not in toml store: %v", err)
	}
}
