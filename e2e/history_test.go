package e2e

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestHistoryShowsCommits(t *testing.T) {
	bin := buildBinary(t)

	// Add first version
	cmd := exec.Command(bin, "add", "hist_entry")
	cmd.Stdin = strings.NewReader("v1\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("add v1 failed: %v\n%s", err, out)
	}

	// Add second version
	cmd = exec.Command(bin, "add", "hist_entry")
	cmd.Stdin = strings.NewReader("v2\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("add v2 failed: %v\n%s", err, out)
	}

	// History should show two commits
	hist := exec.Command(bin, "history", "hist_entry")
	out, err := hist.CombinedOutput()
	if err != nil {
		t.Fatalf("history failed: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "feat: add hist_entry.md") {
		t.Fatalf("expected history to include commit subject, got: %s", out)
	}
}

func TestHistoryMissingEntry(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "history", "missing_entry")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected history to fail for missing entry")
	}
	if !strings.Contains(string(out), "history failed: not found") {
		t.Fatalf("expected not found message, got: %s", out)
	}
}

func TestHistoryReverseLimit(t *testing.T) {
	bin := buildBinary(t)

	// Create three versions
	for i := 0; i < 3; i++ {
		cmd := exec.Command(bin, "add", "hist_limit")
		cmd.Stdin = strings.NewReader("v" + fmt.Sprint(i) + "\n")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("add v%d failed: %v\n%s", i, err, out)
		}
	}

	hist := exec.Command(bin, "history", "hist_limit", "--limit", "2", "--reverse")
	out, err := hist.CombinedOutput()
	if err != nil {
		t.Fatalf("history failed: %v\n%s", err, out)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines from limit, got %d: %v", len(lines), lines)
	}
}
