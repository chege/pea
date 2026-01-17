package cmd

import (
	"os"
	"pea/internal/app"
	"pea/platform"
)

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != 0
}

func copyToClipboard(s string) error {
	// Use platform clipboard abstraction
	if err := platform.ClipboardImpl.Init(); err != nil {
		return err
	}

	return platform.ClipboardImpl.WriteText(s)
}

// Helper to access internal read functionality if needed by other commands (like root)
func ReadEntry(store, name, rev string) ([]byte, error) {
	return app.ReadEntry(store, name, rev)
}
