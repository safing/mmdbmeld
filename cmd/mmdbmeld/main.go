package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/safing/mmdbmeld"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <config.yml>\n", os.Args[0])
		os.Exit(1)
	}

	c, err := mmdbmeld.LoadConfig(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for _, db := range c.Databases {
		fmt.Printf("\n==========\nbuilding %s\n", db.Name)

		sources, err := mmdbmeld.LoadSources(db)
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}

		// Create wait group.
		var wg sync.WaitGroup
		wg.Add(1)

		// Start log writer.
		updates := make(chan string, 100)
		go func() {
			defer wg.Done()
			for msg := range updates {
				fmt.Println(msg)
			}
		}()

		err = mmdbmeld.WriteMMDB(db, sources, updates)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}

		wg.Wait()
	}
}
