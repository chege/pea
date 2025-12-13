package e2e

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompletionBash(t *testing.T) {
	root := filepath.Join("..")
	bin := filepath.Join(root, "bin", "pea")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	out, err := exec.Command(bin, "completion", "bash").CombinedOutput()
	if err != nil {
		t.Fatalf("completion failed: %v\n%s", err, out)
	}
	if len(out) == 0 {
		t.Fatalf("empty completion output")
	}
	if !strings.Contains(string(out), "__start_p") {
		t.Fatalf("completion script missing expected function")
	}
}
