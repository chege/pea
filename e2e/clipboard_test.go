package e2e

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// This test verifies that when output is piped (non-TTY), clipboard is not modified.
func TestClipboardNotModifiedWhenPiped(t *testing.T) {
	root := filepath.Join("..")
	bin := filepath.Join(root, "bin", "pea")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	// Add entry
	add := exec.Command(bin, "add", "cliptest")
	add.Stdin = strings.NewReader("clipboard data\n")
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}
	// Set baseline clipboard
	if out, err := exec.Command("bash", "-c", "printf baseline | pbcopy").CombinedOutput(); err != nil {
		t.Fatalf("pbcopy baseline failed: %v\n%s", err, out)
	}
	// Pipe output to non-TTY and ensure clipboard unchanged
	sh := exec.Command("bash", "-c", bin+" cliptest | cat > /dev/null")
	if out, err := sh.CombinedOutput(); err != nil {
		t.Fatalf("retrieve piped failed: %v\n%s", err, out)
	}
	clip, err := exec.Command("bash", "-c", "pbpaste").CombinedOutput()
	if err != nil {
		t.Fatalf("pbpaste failed: %v", err)
	}
	if string(clip) != "baseline" {
		t.Fatalf("clipboard should remain baseline, got: %q", string(clip))
	}
}
