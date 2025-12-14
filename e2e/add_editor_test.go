package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestAddOpensEditorWhenNoStdin(t *testing.T) {

	bin := buildBinary(t)

	// Create a temp editor script that writes to the file
	tmp := t.TempDir()

	ed := filepath.Join(tmp, "editor.sh")

	script := []byte("#!/bin/bash\nset -eu\nfile=\"$1\"\necho 'edited' > \"$file\"\n")

	if err := os.WriteFile(ed, script, 0o755); err != nil {

		t.Fatal(err)
	}

	cmd := exec.Command(bin, "add", "edited_entry")

	cmd.Env = append(os.Environ(), "EDITOR="+ed)

	out, err := cmd.CombinedOutput()

	if err != nil {

		t.Fatalf("add via editor failed: %v\n%s", err, out)
	}

	home, _ := os.UserHomeDir()

	store := filepath.Join(home, ".pea", "prompts", "edited_entry.md")

	b, err := os.ReadFile(store)

	if err != nil {

		t.Fatalf("missing entry: %v", err)
	}

	if string(b) != "edited\n" {

		t.Fatalf("unexpected content: %q", string(b))
	}
}
