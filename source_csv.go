package mmdbmeld

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

// CSVSource reads geoip data in csv format.
type CSVSource struct {
	file   string
	reader *csv.Reader
	fields []string
	types  map[string]string

	err error
}

// LoadCSVSource returns a new CSVSource.
func LoadCSVSource(input DatabaseInput, types map[string]string) (*CSVSource, error) {
	file, err := os.Open(input.File)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = len(input.Fields)

	return &CSVSource{
		file:   input.File,
		reader: reader,
		fields: input.Fields,
		types:  types,
	}, nil
}

// Name returns an identifying name for the source.
func (csv *CSVSource) Name() string {
	return csv.file
}

// NextEntry returns the next entry of the source.
// If nil, nil is returned, stop reading and check Err().
func (csv *CSVSource) NextEntry() (*SourceEntry, error) {
	// Check if there is an error, do not read if there is an error.
	if csv.err != nil {
		return nil, nil //nolint:nilerr
	}

	// Read and parse line.
	row, err := csv.reader.Read()
	if err != nil {
		csv.err = err
		return nil, nil //nolint:nilerr
	}
	se := &SourceEntry{
		Values: make(map[string]SourceValue),
	}
	for i := 0; i < len(csv.fields); i++ {
		fieldName := csv.fields[i]

		switch fieldName {
		case "from":
			fromIP := net.ParseIP(row[i])
			if fromIP == nil {
				return nil, fmt.Errorf("failed to parse IP %q", row[i])
			}
			se.From = fromIP
			// Force IPv4 representation for IPv4 for better further processing.
			if v4 := se.From.To4(); v4 != nil {
				se.From = v4
			}
		case "to":
			toIP := net.ParseIP(row[i])
			if toIP == nil {
				return nil, fmt.Errorf("failed to parse IP %q", row[i])
			}
			se.To = toIP
			// Force IPv4 representation for IPv4 for better further processing.
			if v4 := se.To.To4(); v4 != nil {
				se.To = v4
			}
		case "", "-":
			// Ignore
		default:
			fieldType, ok := csv.types[fieldName]
			if ok && fieldType != "" && fieldType != "-" && row[i] != "" {
				se.Values[fieldName] = SourceValue{
					Type:  csv.types[fieldName],
					Value: row[i],
				}
			}
		}
	}

	return se, nil
}

// Err returns the processing error encountered by the source.
func (csv *CSVSource) Err() error {
	switch {
	case csv.err == nil:
		return nil
	case errors.Is(csv.err, io.EOF):
		return nil
	default:
		return csv.err
	}
}
