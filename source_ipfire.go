package geoipbuilder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"os"
	"strings"
)

// IPFireSource reads geoip data in the ipfire format.
type IPFireSource struct {
	file     string
	reader   *textproto.Reader
	fieldMap map[string]string
	types    map[string]string

	asOrgCache map[string]string

	err error
}

// LoadIPFireSource returns a new IPFireSource.
func LoadIPFireSource(input DatabaseInput, types map[string]string) (*IPFireSource, error) {
	file, err := os.Open(input.File)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	reader := textproto.NewReader(bufio.NewReader(file))

	// Skip comment section.
	for {
		line, err := reader.ReadLine()
		if err != nil {
			continue
		}
		// Discard lines until the comments stop.
		if !strings.HasPrefix(line, "#") {
			break
		}
	}

	return &IPFireSource{
		file:       input.File,
		reader:     reader,
		fieldMap:   input.FieldMap,
		types:      types,
		asOrgCache: make(map[string]string),
	}, nil
}

// Name returns an identifying name for the source.
func (ipf *IPFireSource) Name() string {
	return ipf.file
}

// NextEntry returns the next entry of the source.
// If nil, nil is returned, stop reading and check Err().
func (ipf *IPFireSource) NextEntry() (*SourceEntry, error) {
	// Check if there is an error, do not read if there is an error.
	if ipf.err != nil {
		return nil, nil //nolint:nilerr
	}

	// Read data sections.
	for {
		// Read next section.
		data, err := ipf.reader.ReadMIMEHeader()
		if err != nil {
			ipf.err = err
			return nil, nil //nolint:nilerr
		}
		// If the section is empty, continue to next.
		if len(data) == 0 {
			continue
		}

		// Parse data.
		se, err := ipf.SourceEntryFromMimeHeader(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ipfire entry: %w", err)
		}

		asNum, ok1 := se.Values[ipf.fieldMap["aut-num"]]
		asOrg, ok2 := se.Values[ipf.fieldMap["name"]]
		// Check for ASN->ASOrg mapping.
		if se.Net == nil {
			if ok1 && ok2 {
				ipf.asOrgCache[asNum.Value] = asOrg.Value
			}
			continue
		}
		// Fill in ASOrg if missing.
		if ok1 && !ok2 {
			asOrgKey := ipf.fieldMap["name"]
			if asOrgKey != "" {
				asOrg := ipf.asOrgCache[asNum.Value]
				if asOrg != "" {
					se.Values[asOrgKey] = SourceValue{
						Type:  ipf.types[asOrgKey],
						Value: asOrg,
					}
				}
			}
		}

		return se, nil
	}
}

// Err returns the processing error encountered by the source.
func (ipf *IPFireSource) Err() error {
	switch {
	case ipf.err == nil:
		return nil
	case errors.Is(ipf.err, io.EOF):
		return nil
	default:
		return ipf.err
	}
}

// SourceEntryFromMimeHeader parses an IPFire entry.
// The keys in the IPFire file are:
// - net
// - aut-num
// - name
// - country
// - is-anycast
// - is-satellite-provider
// - is-anonymous-proxy
// - drop
// .
func (ipf *IPFireSource) SourceEntryFromMimeHeader(data textproto.MIMEHeader) (*SourceEntry, error) {
	se := &SourceEntry{
		Values: make(map[string]SourceValue),
	}

	// Parse Network.
	if netData := data.Get("net"); netData != "" {
		_, ipNet, err := net.ParseCIDR(netData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse net %s: %w", netData, err)
		}
		se.Net = ipNet
	}

	// Parse AS.
	if fieldName, ok := ipf.fieldMap["aut-num"]; ok {
		if asNum := data.Get("aut-num"); asNum != "" {
			asNum := strings.TrimPrefix(asNum, "AS")
			se.Values[fieldName] = SourceValue{
				Type:  ipf.types[fieldName],
				Value: asNum,
			}
		}
	}

	// Parse AS Org.
	if fieldName, ok := ipf.fieldMap["name"]; ok {
		if asOrg := data.Get("name"); asOrg != "" {
			se.Values[fieldName] = SourceValue{
				Type:  ipf.types[fieldName],
				Value: asOrg,
			}
		}
	}

	// Parse country.
	if fieldName, ok := ipf.fieldMap["country"]; ok {
		if country := data.Get("country"); len(country) == 2 {
			se.Values[fieldName] = SourceValue{
				Type:  ipf.types[fieldName],
				Value: country,
			}
		}
	}

	// Parse flags.
	if fieldName, ok := ipf.fieldMap["is-anycast"]; ok {
		if flag := data.Get("is-anycast"); flag != "" {
			se.Values[fieldName] = SourceValue{
				Type:  ipf.types[fieldName],
				Value: "true",
			}
		}
	}
	if fieldName, ok := ipf.fieldMap["is-satellite-provider"]; ok {
		if flag := data.Get("is-satellite-provider"); flag != "" {
			se.Values[fieldName] = SourceValue{
				Type:  ipf.types[fieldName],
				Value: "true",
			}
		}
	}
	if fieldName, ok := ipf.fieldMap["is-anonymous-proxy"]; ok {
		if flag := data.Get("is-anonymous-proxy"); flag != "" {
			se.Values[fieldName] = SourceValue{
				Type:  ipf.types[fieldName],
				Value: "true",
			}
		}
	}
	if fieldName, ok := ipf.fieldMap["drop"]; ok {
		if flag := data.Get("drop"); flag != "" {
			se.Values[fieldName] = SourceValue{
				Type:  ipf.types[fieldName],
				Value: "true",
			}
		}
	}

	return se, nil
}
