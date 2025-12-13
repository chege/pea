package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCompletionBash(t *testing.T) {

	bin := buildBinary(t)

	out, err := exec.Command(bin, "completion", "bash").CombinedOutput()

	if err != nil {

		t.Fatalf("completion failed: %v\n%s", err, out)
	}

	if len(out) == 0 {

		t.Fatalf("empty completion output")
	}

	if !strings.Contains(string(out), "__start_p") {

		t.Fatalf("completion script missing expected function")
	}
}
