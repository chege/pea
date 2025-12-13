package main

import "github.com/spf13/cobra"

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
		if toComplete == "" || hasPrefix(n, toComplete) {
			out = append(out, n)
		}
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

func hasPrefix(s, prefix string) bool {
	if len(prefix) == 0 {
		return true
	}
	if len(prefix) > len(s) {
		return false
	}
	return s[:len(prefix)] == prefix
}
