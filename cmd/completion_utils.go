package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"pea/internal/app"
)

func completeNames(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	store, err := app.EnsureStore()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	names, err := app.ListEntries(store)
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
