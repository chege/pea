package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRetrieveStripsFrontMatter(t *testing.T) {

	bin := buildBinary(t)

	home, _ := os.UserHomeDir()

	store := filepath.Join(home, ".pea", "prompts")

	_ = os.MkdirAll(store, 0o755)

	path := filepath.Join(store, "fm_test.md")

	content := "---\ndescription: test\n---\nBody line 1\nBody line 2\n"

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {

		t.Fatal(err)
	}

	// Retrieve (should strip front matter)
	cmd := exec.Command(bin, "get", "fm_test")
	out, err := cmd.CombinedOutput()
	if err != nil {

		t.Fatalf("retrieve failed: %v\n%s", err, out)
	}

	if string(out) != "Body line 1\nBody line 2\n" {

		t.Fatalf("unexpected body: %q", string(out))
	}
}
