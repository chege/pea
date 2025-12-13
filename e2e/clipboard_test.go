// Package e2e contains end-to-end tests for the application.
package e2e

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// This test verifies that when output is piped (non-TTY), clipboard is not modified.
func TestClipboardNotModifiedWhenPiped(t *testing.T) {
	bin := buildBinary(t)
	// Add entry
	add := exec.Command(bin, "add", "cliptest")
	add.Stdin = strings.NewReader("clipboard data\n")
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}
	// Set baseline clipboard (use fake file when running headless)
	if fake := os.Getenv("PEA_FAKE_CLIP_FILE"); fake != "" {
		if err := os.WriteFile(fake, []byte("baseline"), 0o644); err != nil {
			t.Fatalf("write fake baseline failed: %v", err)
		}
	} else {
		if out, err := exec.Command("bash", "-c", "printf baseline | pbcopy").CombinedOutput(); err != nil {
			t.Fatalf("pbcopy baseline failed: %v\n%s", err, out)
		}
	}
	// Pipe output to non-TTY and ensure clipboard unchanged
	sh := exec.Command("bash", "-c", bin+" cliptest | cat > /dev/null")
	if out, err := sh.CombinedOutput(); err != nil {
		t.Fatalf("retrieve piped failed: %v\n%s", err, out)
	}
	if fake := os.Getenv("PEA_FAKE_CLIP_FILE"); fake != "" {
		clip, err := os.ReadFile(fake)
		if err != nil {
			t.Fatalf("read fake clipboard failed: %v", err)
		}
		if string(clip) != "baseline" {
			t.Fatalf("clipboard should remain baseline, got: %q", string(clip))
		}
	} else {
		clip, err := exec.Command("bash", "-c", "pbpaste").CombinedOutput()
		if err != nil {
			t.Fatalf("pbpaste failed: %v", err)
		}
		if string(clip) != "baseline" {
			t.Fatalf("clipboard should remain baseline, got: %q", string(clip))
		}
	}
}
