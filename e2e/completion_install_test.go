package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompletionInstallUpdatesWithMarkers(t *testing.T) {
	bin := buildBinary(t)
	home := t.TempDir()

	// 1. First Install
	cmd := exec.Command(bin, "completion", "install")
	cmd.Env = append(os.Environ(), "HOME="+home, "SHELL=/bin/bash")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("install failed: %v\n%s", err, out)
	}

	bashrc := filepath.Join(home, ".bashrc")
	content, err := os.ReadFile(bashrc)
	if err != nil {
		t.Fatal(err)
	}
	s := string(content)

	startMarker := "# PEA_COMPLETION_START"
	endMarker := "# PEA_COMPLETION_END"

	if !strings.Contains(s, startMarker) || !strings.Contains(s, endMarker) {
		t.Fatalf("markers missing in first install:\n%s", s)
	}

	// 2. Modify the file manually (simulate user change or drift, or verify replacement)
	// We'll corrupt the content between markers to ensure it gets fixed/replaced

	brokenContent := strings.Replace(s, "source ", "# broken ", 1)
	if err := os.WriteFile(bashrc, []byte(brokenContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// 3. Second Install (Should repair/update)
	cmd2 := exec.Command(bin, "completion", "install")
	cmd2.Env = append(os.Environ(), "HOME="+home, "SHELL=/bin/bash")
	if out, err := cmd2.CombinedOutput(); err != nil {
		t.Fatalf("re-install failed: %v\n%s", err, out)
	}

	content2, err := os.ReadFile(bashrc)
	if err != nil {
		t.Fatal(err)
	}
	s2 := string(content2)

	if strings.Contains(s2, "# broken") {
		t.Fatalf("re-install failed to replace content:\n%s", s2)
	}
	if !strings.Contains(s2, "source ") {
		t.Fatalf("re-install failed to restore source command:\n%s", s2)
	}

	// Ensure markers are still there
	if !strings.Contains(s2, startMarker) || !strings.Contains(s2, endMarker) {
		t.Fatalf("markers missing after update:\n%s", s2)
	}

	// Ensure no duplicate blocks
	if strings.Count(s2, startMarker) > 1 {
		t.Fatalf("duplicate blocks detected:\n%s", s2)
	}
}
