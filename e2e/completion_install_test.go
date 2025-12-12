package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCompletionInstall(t *testing.T) {
	root := filepath.Join("..")
	bin := filepath.Join(root, "bin", "pea")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	// Install
	out, err := exec.Command(bin, "completion", "install").CombinedOutput()
	if err != nil {
		t.Fatalf("completion install failed: %v\n%s", err, out)
	}
	// Verify files exist
	home, _ := os.UserHomeDir()
	base := filepath.Join(home, ".pea")
	if _, err := os.Stat(filepath.Join(base, "p.bash")); err != nil {
		t.Fatalf("bash completion missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(base, "_p")); err != nil {
		t.Fatalf("zsh completion missing: %v", err)
	}
}
