package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

func TestListEntries(t *testing.T) {
	bin := buildBinary(t)
	// Create entries
	add1 := exec.Command(bin, "add", "b_entry")
	add1.Stdin = strings.NewReader("b\n")
	if out, err := add1.CombinedOutput(); err != nil {
		t.Fatalf("add1 failed: %v\n%s", err, out)
	}
	add2 := exec.Command(bin, "add", "a_entry")
	add2.Stdin = strings.NewReader("a\n")
	if out, err := add2.CombinedOutput(); err != nil {
		t.Fatalf("add2 failed: %v\n%s", err, out)
	}
	// List
	out, err := exec.Command(bin, "ls").CombinedOutput()
	if err != nil {
		t.Fatalf("ls failed: %v\n%s", err, out)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	// Must contain at least the two entries in sorted order as a subsequence
	want := []string{"a_entry", "b_entry"}
	idx := 0
	for _, l := range lines {
		if idx < len(want) && l == want[idx] {
			idx++
		}
	}
	if idx != len(want) {
		t.Fatalf("ls output missing ordered entries: got %q", string(out))
	}
}
