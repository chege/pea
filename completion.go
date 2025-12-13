package main

import (
	"strings"

	"github.com/spf13/cobra"
)

func completeNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	store, err := ensureStore()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	names, err := listEntries(store)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var out []string
	for _, n := range names {
		if strings.HasPrefix(n, toComplete) || toComplete == "" {
			out = append(out, n)
		}
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}
