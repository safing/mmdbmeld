package geoipbuilder

import (
	"fmt"
	"testing"
)

func TestMMDBTypes(t *testing.T) {
	t.Parallel()

	entry := &SourceEntry{
		Values: map[string]SourceValue{
			"country.iso_code": {
				Type:  "string",
				Value: "AT",
			},
			"location.latitude": {
				Type:  "float32",
				Value: "48.1230",
			},
			"location.longitude": {
				Type:  "float32",
				Value: "16.2221",
			},
			"autonomous_system_organization": {
				Type:  "string",
				Value: "Example Org",
			},
			"autonomous_system_number": {
				Type:  "uint32",
				Value: "12345",
			},
		},
	}
	mmdbMapString := "map[autonomous_system_number:12345 autonomous_system_organization:Example Org country:map[iso_code:AT] location:map[latitude:48.12 longitude:16.22]]"

	mmdbMap, err := entry.ToMMDBMap(Optimizations{
		FloatDecimals: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	s := fmt.Sprintf("%v", mmdbMap)
	if s != mmdbMapString {
		t.Fatalf("mmdb map string not as expected: %s", s)
	}
}
