package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenameEntry(t *testing.T) {

	bin := buildBinary(t)

	// Create entry
	add := exec.Command(bin, "add", "alpha")

	add.Stdin = strings.NewReader("A\n")

	if out, err := add.CombinedOutput(); err != nil {

		t.Fatalf("add failed: %v\n%s", err, out)
	}

	// Rename
	mv := exec.Command(bin, "mv", "alpha", "beta")

	if out, err := mv.CombinedOutput(); err != nil {

		t.Fatalf("mv failed: %v\n%s", err, out)
	}

	// Verify new file content
	home, _ := os.UserHomeDir()

	path := filepath.Join(home, ".pea", "prompts", "beta.md")

	b, err := os.ReadFile(path)

	if err != nil {

		t.Fatalf("beta missing: %v", err)
	}

	if string(b) != "A\n" {

		t.Fatalf("unexpected content: %q", string(b))
	}
}
