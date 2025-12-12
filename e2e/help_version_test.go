package e2e

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	root := filepath.Join("..")
	bin := filepath.Join(root, "bin", "pea")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	out, err := exec.Command(bin, "--version").CombinedOutput()
	if err != nil {
		t.Fatalf("version failed: %v\n%s", err, out)
	}
	if !strings.HasPrefix(string(out), "p version ") {
		t.Fatalf("unexpected version output: %q", string(out))
	}
}

func TestHelp(t *testing.T) {
	root := filepath.Join("..")
	bin := filepath.Join(root, "bin", "pea")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	out, err := exec.Command(bin).CombinedOutput()
	if err != nil {
		t.Fatalf("running without args failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "p is a fast, local CLI") {
		t.Fatalf("help should contain description, got: %q", string(out))
	}
}
