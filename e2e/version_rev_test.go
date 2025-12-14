package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

func TestGetAtRevision(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "add", "rev_entry")
	cmd.Stdin = strings.NewReader("v1\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("add v1 failed: %v\n%s", err, out)
	}

	cmd = exec.Command(bin, "add", "rev_entry")
	cmd.Stdin = strings.NewReader("v2\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("add v2 failed: %v\n%s", err, out)
	}

	get := exec.Command(bin, "--rev", "HEAD~1", "rev_entry")
	out, err := get.CombinedOutput()
	if err != nil {
		t.Fatalf("get at rev failed: %v\n%s", err, out)
	}
	if string(out) != "v1\n" {
		t.Fatalf("expected v1, got: %q", string(out))
	}
}

func TestGetAtMissingRevision(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "--rev", "doesnotexist", "missing_rev")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected failure for missing ref")
	}
	if !strings.Contains(string(out), "not found in ref") {
		t.Fatalf("expected not found in ref message, got: %s", out)
	}
}
