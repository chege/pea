package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestAddFallsBackToBrowserWhenEditorUnset(t *testing.T) {
	bin := buildBinary(t)

	tmp := t.TempDir()
	browserScript := filepath.Join(tmp, "browser.sh")
	script := []byte("#!/bin/bash\nset -eu\nurl=\"$1\"\npath=${url#file://}\nprintf 'fallback\\n' > \"$path\"\n")
	if err := os.WriteFile(browserScript, script, 0o755); err != nil {
		t.Fatal(err)
	}

	store := filepath.Join(tmp, "store")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(bin, "add", "fallback_entry")
	cmd.Env = append(os.Environ(), "PEA_STORE="+store, "EDITOR=", "BROWSER="+browserScript)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("add via browser fallback failed: %v\n%s", err, out)
	}

	entry := filepath.Join(store, "fallback_entry.txt")
	data, err := os.ReadFile(entry)
	if err != nil {
		t.Fatalf("missing entry: %v", err)
	}
	if string(data) != "fallback\n" {
		t.Fatalf("unexpected content: %q", string(data))
	}
	_ = os.Remove(entry)
}
