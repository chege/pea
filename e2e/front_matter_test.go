package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRetrieveStripsFrontMatter(t *testing.T) {
	root := filepath.Join("..")
	bin := filepath.Join(root, "bin", "pea")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	home, _ := os.UserHomeDir()
	store := filepath.Join(home, ".pea", "prompts")
	_ = os.MkdirAll(store, 0o755)
	path := filepath.Join(store, "fm_test.txt")
	content := "---\ndescription: test\n---\nBody line 1\nBody line 2\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command(bin, "fm_test").CombinedOutput()
	if err != nil {
		t.Fatalf("retrieve failed: %v\n%s", err, out)
	}
	if string(out) != "Body line 1\nBody line 2\n" {
		t.Fatalf("unexpected body: %q", string(out))
	}
}
