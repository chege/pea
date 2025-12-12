package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnvOverridesStore(t *testing.T) {
	root := filepath.Join("..")
	bin := filepath.Join(root, "bin", "pea")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	tmp := t.TempDir()
	cmd := exec.Command(bin, "add", "envcase")
	cmd.Env = append(os.Environ(), "PEA_STORE="+tmp)
	cmd.Stdin = strings.NewReader("x\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(tmp, "envcase.txt")); err != nil {
		t.Fatalf("file not in env store: %v", err)
	}
}
