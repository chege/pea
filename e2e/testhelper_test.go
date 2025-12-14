// Package e2e contains end-to-end tests for the application.
package e2e

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
)

var (
	binPath   string
	buildOnce sync.Once
)

func buildBinary(t *testing.T) string {

	t.Helper()

	buildOnce.Do(func() {

		root, err := findModuleRoot()

		if err != nil {
			t.Fatalf("find module root: %v", err)
		}

		tmpDir, err := os.MkdirTemp("", "pea-bin-")

		if err != nil {
			t.Fatalf("mkdtemp failed: %v", err)
		}

		bin := filepath.Join(tmpDir, "pea")

		cmd := exec.Command("go", "build", "-o", bin, ".")

		cmd.Dir = root

		if out, err := cmd.CombinedOutput(); err != nil {

			t.Fatalf("build failed: %v\n%s", err, out)
		}

		// Create a small wrapper that enables headless/fake platform implementations
		// for tests (PEA_HEADLESS) and points the fake clipboard file to a location
		// under the same tmp dir so tests can inspect it.
		wrapper := filepath.Join(tmpDir, "pea-wrapper")

		fakeClip := filepath.Join(tmpDir, "pea_fake_clipboard")
		home := filepath.Join(tmpDir, "home")
		_ = os.MkdirAll(home, 0o755)

		script := "#!/bin/bash\nif [ -z \"$HOME\" ]; then export HOME='" + home + "'; fi\nexport PEA_HEADLESS=1\nexport PEA_FAKE_CLIP_FILE='" + fakeClip + "'\nexec '" + bin + "' \"$@\"\n"

		// Also set the environment in the test process so tests can use the same
		// fake clipboard path without falling back to OS clipboard utilities.
		_ = os.Setenv("HOME", home)
		_ = os.Setenv("PEA_HEADLESS", "1")
		_ = os.Setenv("PEA_FAKE_CLIP_FILE", fakeClip)

		if err := os.WriteFile(wrapper, []byte(script), 0o755); err != nil {

			t.Fatalf("write wrapper failed: %v", err)
		}

		binPath = wrapper
	})

	return binPath
}

func findModuleRoot() (string, error) {

	dir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	for {

		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {

			return dir, nil
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			return "", errors.New("go.mod not found")
		}

		dir = parent
	}
}
