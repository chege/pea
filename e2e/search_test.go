package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	bin := buildBinary(t)

	// Create entries with tags
	entries := []struct {
		name    string
		content string
	}{
		{
			name: "golang_basics",
			content: `---
tags: [coding, go]
---
Basics of Go.`,
		},
		{
			name: "rust_basics",
			content: `---
tags: [coding, rust]
---
Basics of Rust.`,
		},
		{
			name: "cooking_pasta",
			content: `---
tags: [hobby, food]
---
How to cook pasta.`,
		},
	}

	for _, e := range entries {
		cmd := exec.Command(bin, "add", e.name)
		cmd.Stdin = strings.NewReader(e.content)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to add %s: %v\n%s", e.name, err, out)
		}
	}

	tests := []struct {
		name     string
		args     []string
		want     []string
		notWant  []string
	}{
		{
			name: "search by substring",
			args: []string{"search", "basics"},
			want: []string{"golang_basics", "rust_basics"},
			notWant: []string{"cooking_pasta"},
		},
		{
			name: "search by tag",
			args: []string{"search", "--tag", "go"},
			want: []string{"golang_basics"},
			notWant: []string{"rust_basics", "cooking_pasta"},
		},
		{
			name: "search by multiple tags",
			args: []string{"search", "--tag", "coding", "--tag", "rust"},
			want: []string{"rust_basics"},
			notWant: []string{"golang_basics", "cooking_pasta"},
		},
		{
			name: "search by substring and tag",
			args: []string{"search", "pasta", "--tag", "food"},
			want: []string{"cooking_pasta"},
			notWant: []string{"golang_basics", "rust_basics"},
		},
		{
			name: "search no match",
			args: []string{"search", "nonexistent"},
			want: []string{}, // empty slice
			notWant: []string{"golang_basics", "rust_basics", "cooking_pasta"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := exec.Command(bin, tt.args...).CombinedOutput()
			if err != nil {
				t.Fatalf("search failed: %v\n%s", err, out)
			}
			got := string(out)
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("expected output to contain %q, but got %q", w, got)
				}
			}
			for _, nw := range tt.notWant {
				if strings.Contains(got, nw) {
					t.Errorf("expected output NOT to contain %q, but got %q", nw, got)
				}
			}
		})
	}
}
