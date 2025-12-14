package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func storePath() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".pea", "prompts")
}

func TestRmDryRunDoesNotDelete(t *testing.T) {
	bin := buildBinary(t)
	store := storePath()

	cmd := exec.Command(bin, "add", "safe_entry")
	cmd.Stdin = strings.NewReader("content\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}

	dry := exec.Command(bin, "rm", "safe_entry", "--dry-run")
	if out, err := dry.CombinedOutput(); err != nil {
		t.Fatalf("dry-run rm failed: %v\n%s", err, out)
	}

	if _, err := os.Stat(filepath.Join(store, "safe_entry.md")); err != nil {
		t.Fatalf("entry should still exist after dry-run: %v", err)
	}
}

func TestRmUndoRestores(t *testing.T) {
	bin := buildBinary(t)
	store := storePath()

	cmd := exec.Command(bin, "add", "undo_entry")
	cmd.Stdin = strings.NewReader("content\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}

	rm := exec.Command(bin, "rm", "undo_entry")
	if out, err := rm.CombinedOutput(); err != nil {
		t.Fatalf("rm failed: %v\n%s", err, out)
	}

	if _, err := os.Stat(filepath.Join(store, "undo_entry.md")); !os.IsNotExist(err) {
		t.Fatalf("entry should be deleted, got err=%v", err)
	}

	undo := exec.Command(bin, "rm", "undo_entry", "--undo")
	if out, err := undo.CombinedOutput(); err != nil {
		t.Fatalf("undo failed: %v\n%s", err, out)
	}

	if _, err := os.Stat(filepath.Join(store, "undo_entry.md")); err != nil {
		t.Fatalf("entry should be restored: %v", err)
	}
}

func TestMvDryRunDoesNotRename(t *testing.T) {
	bin := buildBinary(t)
	store := storePath()

	cmd := exec.Command(bin, "add", "mv_entry")
	cmd.Stdin = strings.NewReader("content\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("add failed: %v\n%s", err, out)
	}

	mv := exec.Command(bin, "mv", "mv_entry", "mv_entry_new", "--dry-run")
	if out, err := mv.CombinedOutput(); err != nil {
		t.Fatalf("mv dry-run failed: %v\n%s", err, out)
	}

	if _, err := os.Stat(filepath.Join(store, "mv_entry.md")); err != nil {
		t.Fatalf("original should remain after dry-run: %v", err)
	}
	if _, err := os.Stat(filepath.Join(store, "mv_entry_new.md")); err == nil {
		t.Fatalf("new file should not exist after dry-run")
	}
}
