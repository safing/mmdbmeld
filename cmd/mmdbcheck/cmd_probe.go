package main

import (
	"fmt"
	"net"
	"runtime"
	"sync"

	"github.com/oschwald/maxminddb-golang"
	"github.com/spf13/cobra"
)

var probeCommand = &cobra.Command{
	Use:  "probe",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reader, err := openMMDB(args[0])
		if err != nil {
			return err
		}
		defer reader.Close() //nolint:errcheck

		probe(reader)
		return nil
	},
}

func probe(reader *maxminddb.Reader) {
	var lookupErrors int
	var msLock sync.Mutex
	ms := &mmdbStats{}

	workers := runtime.NumCPU() * 2
	workFeed := make(chan net.IP, workers*100)
	var wg sync.WaitGroup
	wg.Add(workers + 1)

	switch reader.Metadata.IPVersion {
	case 4:
		go func() {
			var a, b, c uint16
			for a = 1; a <= 254; a++ { // Skip 1. and 255.
				for b = 0; b <= 255; b++ {
					for c = 0; c <= 255; c++ {
						workFeed <- net.IPv4(byte(a), byte(b), byte(c), byte((a*b*c)%256))
					}
				}
				fmt.Print(".")
			}
			fmt.Println("::")
			close(workFeed)
			wg.Done()
		}()
	case 6:
		go func() {
			// Limit search to 2000::/4 for now, as this is the actually used global space.
			fmt.Println("Probing 2000::/4 in /32 blocks")
			var a, b, c, d uint16
			for a = 0x20; a <= 0x30; a++ {
				for b = 0; b <= 255; b++ {
					for c = 0; c <= 255; c++ {
						for d = 0; d <= 255; d++ {
							workFeed <- net.IP([]byte{
								byte(a),
								byte(b),
								byte(c),
								byte(d),
								byte((a * a) % 256),
								byte((a * b) % 256),
								byte((a * c) % 256),
								byte((a * d) % 256),
								byte((b * b) % 256),
								byte((b * c) % 256),
								byte((b * d) % 256),
								byte((c * c) % 256),
								byte((c * d) % 256),
								byte((a * a * a) % 256),
								byte((a * a * b) % 256),
								byte((a * a * c) % 256),
							})
						}
					}
					fmt.Print(".")
				}
				fmt.Println(":")
			}
			fmt.Println("::")
			close(workFeed)
			wg.Done()
		}()
	default:
		fmt.Printf("unsupported IP version %d for probing\n", reader.Metadata.IPVersion)
		return
	}

	for i := 0; i < workers; i++ {
		go func() {
			for ip := range workFeed {
				// Only check global IP scope.
				if GetIPScope(ip) != Global {
					continue
				}

				l := &Location{}
				err := reader.Lookup(ip, l)
				msLock.Lock()
				if err != nil {
					lookupErrors++
				} else {
					ms.Add(l)
				}
				msLock.Unlock()
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Printf("analyzed with %d lookup errors\n", lookupErrors)
	ms.Print()
}
