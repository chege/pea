package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

func TestDeleteEntry(t *testing.T) {
	bin := buildBinary(t)
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
