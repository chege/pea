package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddRejectsInvalidName(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "add", "!!!")
	cmd.Stdin = strings.NewReader("content\n")
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("expected failure for invalid name, got success: %s", out)
	} else if !strings.Contains(string(out), "invalid name") {
		t.Fatalf("expected invalid name message, got: %s", out)
	}
}

func TestRmNotFoundMessage(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "rm", "does_not_exist")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected rm to fail for missing entry")
	}
	if !strings.Contains(string(out), "delete failed: not found") {
		t.Fatalf("expected not found message, got: %s", out)
	}
}

func TestMvRejectsInvalidNames(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "mv", "good", "!!!")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected mv to fail for invalid new name")
	}
	if !strings.Contains(string(out), "invalid name") {
		t.Fatalf("expected invalid name message, got: %s", out)
	}
}

func TestCompletionInstallErrorsAreContextual(t *testing.T) {
	bin := buildBinary(t)

	tmp := filepath.Join(os.TempDir(), "pea-home-file")
	_ = os.RemoveAll(tmp)
	if err := os.WriteFile(tmp, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup home file: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(tmp) })

	cmd := exec.Command(bin, "completion", "install")
	cmd.Env = append(os.Environ(), "HOME="+tmp)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected completion install to fail when HOME is a file, got: %s", out)
	}
	if !strings.Contains(string(out), "create dir") {
		t.Fatalf("expected contextual error, got: %s", out)
	}
}
