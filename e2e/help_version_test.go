package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {

	bin := buildBinary(t)

	out, err := exec.Command(bin, "--version").CombinedOutput()

	if err != nil {

		t.Fatalf("version failed: %v\n%s", err, out)
	}

	if !strings.HasPrefix(string(out), "pea version ") {

		t.Fatalf("unexpected version output: %q", string(out))
	}
}

func TestHelp(t *testing.T) {

	bin := buildBinary(t)

	out, err := exec.Command(bin).CombinedOutput()

	if err != nil {

		t.Fatalf("running without args failed: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "pea is a fast, local CLI") {

		t.Fatalf("help should contain description, got: %q", string(out))
	}
}
