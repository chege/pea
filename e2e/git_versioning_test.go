package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGitCommitsOnAdd(t *testing.T) {
	root := filepath.Join("..")
	bin := filepath.Join(root, "bin", "pea")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	// Use temp store
	tmp := t.TempDir()
	cmd := exec.Command(bin, "add", "gitunit")
	cmd.Env = append(os.Environ(), "PEA_STORE="+tmp)
	cmd.Stdin = strings.NewReader("v1\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}
	// Check latest commit
	out, err := exec.Command("bash", "-c", "cd '"+tmp+"' && git log --oneline -n 1").CombinedOutput()
	if err != nil {
		t.Fatalf("git log failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "add gitunit") {
		t.Fatalf("expected commit message, got: %q", string(out))
	}
}
