package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddEmptyContentFailsAndCleansUp(t *testing.T) {
	bin := buildBinary(t)
	home := os.Getenv("HOME")
	store := filepath.Join(home, ".pea", "prompts")

	cmd := exec.Command(bin, "add", "empty_case")
	cmd.Stdin = strings.NewReader("   \n\t\n")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected add to fail for empty content")
	}
	if !strings.Contains(string(out), "empty content") {
		t.Fatalf("expected empty content message, got: %s", out)
	}

	if _, err := os.Stat(filepath.Join(store, "empty_case.md")); err == nil {
		t.Fatalf("expected file to be absent after empty content")
	}
}
