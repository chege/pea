package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddViaStdin(t *testing.T) {
	root := filepath.Join("..")
	bin := filepath.Join(root, "bin", "pea")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	cmd := exec.Command(bin, "add", "hello_world")
	cmd.Stdin = strings.NewReader("Hello world\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}
	store, _ := os.UserHomeDir()
	store = filepath.Join(store, ".pea", "prompts", "hello_world.txt")
	data, err := os.ReadFile(store)
	if err != nil {
		t.Fatalf("entry not written: %v", err)
	}
	if string(data) != "Hello world\n" {
		t.Fatalf("unexpected content: %q", string(data))
	}
}
