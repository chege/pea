// Package platform provides adapters and fakes for OS-dependent features
// such as clipboard and browser interactions used by the application. The
// real adapters delegate to OS libraries while the fakes enable headless
// CI and tests.
package platform

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	browserpkg "github.com/pkg/browser"
	clipboardpkg "golang.design/x/clipboard"
)

// Clipboard abstracts clipboard operations used by the application.
type Clipboard interface {
	Init() error
	WriteText(string) error
}

// Browser abstracts opening files/URLs with the OS default handler.
type Browser interface {
	OpenFile(string) error
}

var (
	// ClipboardImpl is the implementation used by the application. It may be
	// swapped for a fake in tests by setting PEA_HEADLESS in the environment.
	ClipboardImpl Clipboard
	// BrowserImpl is the implementation used by the application.
	BrowserImpl Browser
	once        sync.Once
)

func init() {
	once.Do(initImpls)
}

func initImpls() {
	if os.Getenv("PEA_HEADLESS") != "" {
		// In test/headless mode, use fakes that do not require OS services.
		ClipboardImpl = &fakeClipboard{}
		BrowserImpl = &fakeBrowser{}
		return
	}
	// Default real implementations
	ClipboardImpl = &realClipboard{}
	BrowserImpl = &realBrowser{}
}

// realClipboard is an adapter to golang.design/x/clipboard.
type realClipboard struct{}

func (r *realClipboard) Init() error {
	return clipboardpkg.Init()
}

func (r *realClipboard) WriteText(s string) error {
	clipboardpkg.Write(clipboardpkg.FmtText, []byte(s))
	return nil
}

// realBrowser is an adapter to github.com/pkg/browser.
type realBrowser struct{}

func (r *realBrowser) OpenFile(path string) error {
	return browserpkg.OpenFile(path)
}

// fakeClipboard writes clipboard contents to a file so tests can inspect it.
type fakeClipboard struct{}

func (f *fakeClipboard) Init() error { return nil }

func (f *fakeClipboard) WriteText(s string) error {
	p := os.Getenv("PEA_FAKE_CLIP_FILE")
	if p == "" {
		// default to temp file in current working directory
		tmp := os.TempDir()
		p = filepath.Join(tmp, "pea_fake_clipboard")
	}
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, []byte(s), 0o644)
}

// fakeBrowser does a no-op for headless tests.
type fakeBrowser struct{}

func (f *fakeBrowser) OpenFile(_ string) error { return nil }
