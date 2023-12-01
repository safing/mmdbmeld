package main

import "github.com/spf13/cobra"

func getRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "mmdb-check",
		Short: "Check a mmdb database and print results",
	}

	root.AddCommand(
		allChecksCommand,
		masksCommand,
		probeCommand,
	)

	return root
}
