package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

func TestAddRejectsReservedNames(t *testing.T) {
	bin := buildBinary(t)

	reserved := []string{"add", "ls", "rm", "mv", "history", "search", "completion", "help"}

	for _, name := range reserved {
		cmd := exec.Command(bin, "add", name)
		cmd.Stdin = strings.NewReader("content")
		out, err := cmd.CombinedOutput()

		if err == nil {
			t.Errorf("add %s should have failed", name)
		}

		if !strings.Contains(string(out), "reserved command") {
			t.Errorf("add %s error message missing 'reserved command': %s", name, out)
		}
	}
}
