package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/oschwald/maxminddb-golang"
)

func openMMDB(path string) (*maxminddb.Reader, error) {
	// Print some file info.
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	fmt.Printf(
		"loading %s with %.2f MB\n",
		path,
		float64(stat.Size())/1000000,
	)

	// Open gzip in memory.
	if strings.HasSuffix(path, ".gz") {
		// Read db to memory.
		mmdbGzipData, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		// Uncompress db.
		gzipData := bytes.NewBuffer(nil)
		gzipReader, err := gzip.NewReader(bytes.NewBuffer(mmdbGzipData))
		if err != nil {
			return nil, err
		}
		_, err = gzipData.ReadFrom(gzipReader)
		if err != nil {
			return nil, err
		}

		// Print decompression stats.
		fmt.Printf(
			"loaded into memory and decompressed to %.2f MB\n",
			float64(gzipData.Len())/1000000,
		)

		return maxminddb.FromBytes(gzipData.Bytes())
	}

	// Load file directly.
	return maxminddb.Open(path)
}
