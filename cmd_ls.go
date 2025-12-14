package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func addListCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "list stored entries",
		RunE: func(cmd *cobra.Command, _ []string) error {
			store, err := ensureStore()
			if err != nil {
				return err
			}
			entries, err := listEntries(store)
			if err != nil {
				return err
			}
			for _, e := range entries {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), e); err != nil {
					return err
				}
			}
			return nil
		},
	}
	root.AddCommand(cmd)
}

func listEntries(store string) ([]string, error) {
	files, err := os.ReadDir(store)
	if err != nil {
		return nil, err
	}
	nameSet := make(map[string]struct{})
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		switch {
		case strings.HasSuffix(name, defaultExt):
			name = strings.TrimSuffix(name, defaultExt)
			nameSet[name] = struct{}{}
		case strings.HasSuffix(name, legacyExt):
			name = strings.TrimSuffix(name, legacyExt)
			if _, exists := nameSet[name]; !exists {
				nameSet[name] = struct{}{}
			}
		}
	}
	var names []string
	for n := range nameSet {
		names = append(names, n)
	}
	sort.Strings(names)
	return names, nil
}
