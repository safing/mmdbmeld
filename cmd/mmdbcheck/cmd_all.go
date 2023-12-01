package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var allChecksCommand = &cobra.Command{
	Use:  "all",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reader, err := openMMDB(args[0])
		if err != nil {
			return err
		}
		defer reader.Close() //nolint:errcheck

		fmt.Println("\nRunning all checks:")
		fmt.Println("\nProbing:")
		probe(reader)
		fmt.Println("\nNetwork Mask Stats:")
		return masks(reader)
	},
}
