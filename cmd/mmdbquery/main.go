package main

import (
	"fmt"
	"net"
	"os"
	"slices"
	"strings"

	reader "github.com/oschwald/maxminddb-golang"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s <mmdb> <IPs>\n", os.Args[0])
		os.Exit(1)
	}

	r, err := reader.Open(os.Args[1])
	if err != nil {
		fmt.Printf("error:  failed to open mmdb file: %s\n", err)
		os.Exit(2)
	}

	for i, ipArg := range os.Args[2:] {
		// Add a space after first result.
		if i > 0 {
			fmt.Println()
		}

		// Parse IP.
		ip := net.ParseIP(ipArg)
		if ip == nil {
			fmt.Printf("error: invalid IP: %s\n", ipArg)
			os.Exit(2)
		}

		// Get data of IP.
		anyData := make(map[string]any)
		recordNet, ok, err := r.LookupNetwork(ip, &anyData)
		if err != nil {
			fmt.Printf("error: failed to lookup IP: %s\n", err)
			os.Exit(4)
		}
		if !ok {
			fmt.Printf("error: database does not have a record for %s\n", ip)
		}

		// Print result.
		fmt.Printf("%s:\n", recordNet)
		printData(anyData, "  ")
	}
}

type keyValue struct {
	key   string
	value any
}

func printData(m map[string]any, indent string) {
	// Order data.
	ordered := make([]keyValue, 0, len(m))
	for k, v := range m {
		ordered = append(ordered, keyValue{
			key:   k,
			value: v,
		})
	}
	slices.SortFunc[[]keyValue, keyValue](
		ordered,
		func(a, b keyValue) int {
			return strings.Compare(a.key, b.key)
		},
	)

	// Print ordered data.
	for _, e := range ordered {
		subMap, ok := e.value.(map[string]any)
		if ok {
			fmt.Printf("%s%s:\n", indent, e.key)
			printData(subMap, indent+"  ")
		} else {
			fmt.Printf("%s%s: %v\n", indent, e.key, e.value)
		}
	}
}
