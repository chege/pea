package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRetrieveByName(t *testing.T) {

	bin := buildBinary(t)

	// Add entry via stdin
	cmd := exec.Command(bin, "add", "greet")

	cmd.Stdin = strings.NewReader("Hello\n")

	if out, err := cmd.CombinedOutput(); err != nil {

		t.Fatalf("add failed: %v\n%s", err, out)
	}

	// Retrieve
	out, err := exec.Command(bin, "get", "greet").CombinedOutput()
	if err != nil {
		t.Fatalf("retrieve failed: %v\n%s", err, out)

	}

	if string(out) != "Hello\n" {

		t.Fatalf("unexpected output: %q", string(out))
	}

	// Ensure file exists
	home, _ := os.UserHomeDir()

	path := filepath.Join(home, ".pea", "prompts", "greet.md")

	if _, err := os.Stat(path); err != nil {

		t.Fatalf("entry missing: %v", err)
	}
}
