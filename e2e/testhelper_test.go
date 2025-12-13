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
		binPath = bin
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
