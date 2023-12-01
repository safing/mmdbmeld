package main

import (
	"fmt"

	"github.com/oschwald/maxminddb-golang"
	"github.com/spf13/cobra"
)

var masksCommand = &cobra.Command{
	Use:  "masks",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reader, err := openMMDB(args[0])
		if err != nil {
			return err
		}
		defer reader.Close() //nolint:errcheck

		return masks(reader)
	},
}

func masks(reader *maxminddb.Reader) error {
	ms := &mmdbStats{}
	pds := &perDepthStats{}

	iter := reader.Networks()
	for iter.Next() {
		l := &Location{}
		net, err := iter.Network(l)
		if err != nil {
			return fmt.Errorf("failed to get network: %w", err)
		}
		ms.Add(l)
		ones, _ := net.Mask.Size()
		pds.Add(ones, l)
	}

	fmt.Print("          ") // Line up first line with table.
	ms.Print()
	pds.Print()
	return nil
}
