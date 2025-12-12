package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func readEntry(store, name string) ([]byte, error) {
	name = toSnake(name)
	path := filepath.Join(store, name+".txt")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("not found: %s", name)
	}
	return b, nil
}

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil { return false }
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func copyToClipboard(s string) error {
	// macOS: use pbcopy via exec.Command
	c := exec.Command("bash", "-c", "pbcopy")
	in, err := c.StdinPipe()
	if err != nil { return err }
	if err := c.Start(); err != nil { return err }
	if _, err := io.WriteString(in, s); err != nil { return err }
	in.Close()
	return c.Wait()
}
