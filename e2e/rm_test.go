package e2e

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeleteEntry(t *testing.T) {
	root := filepath.Join("..")
	bin := filepath.Join(root, "bin", "pea")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	// Create entry
	add := exec.Command(bin, "add", "temp")
	add.Stdin = strings.NewReader("temp\n")
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}
	// Delete
	rm := exec.Command(bin, "rm", "temp")
	out, err := rm.CombinedOutput()
	if err != nil {
		t.Fatalf("rm failed: %v\n%s", err, out)
	}
	// Ensure retrieval fails
	get := exec.Command(bin, "temp")
	if err := get.Run(); err == nil {
		t.Fatalf("expected retrieve to fail after delete")
	}
}
