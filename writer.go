package mmdbmeld

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"go4.org/netipx"
)

const reportSlotSize = 100_000

// WriteMMDB writes a mmdb file using given config and sources.
// Supply an updates channel to receive update messages about the progress.
func WriteMMDB(dbConfig DatabaseConfig, sources []Source, updates chan string) error {
	// Init writer.
	opts := mmdbwriter.Options{
		DatabaseType:            dbConfig.Name,
		IncludeReservedNetworks: true,
		DisableIPv4Aliasing:     true,
		IPVersion:               dbConfig.MMDB.IPVersion,
		RecordSize:              dbConfig.MMDB.RecordSize,
	}
	writer, err := mmdbwriter.New(opts)
	if err != nil {
		return fmt.Errorf("failed to create mmdb writer for %s: %w", dbConfig.Name, err)
	}
	sendUpdate(updates, fmt.Sprintf(
		"database options set: IPVersion=%d RecordSize=%d (IncludeReservedNetworks=%v DisableIPv4Aliasing=%v)",
		opts.IPVersion,
		opts.RecordSize,
		opts.IncludeReservedNetworks,
		opts.DisableIPv4Aliasing,
	))
	typeKeys := make([]string, 0, len(dbConfig.Types))
	for k, v := range dbConfig.Types {
		if v != "-" && v != "" {
			typeKeys = append(typeKeys, k)
		}
	}
	slices.Sort[[]string, string](typeKeys)
	sendUpdate(updates, fmt.Sprintf(
		"database types: %s",
		strings.Join(typeKeys, ", "),
	))
	sendUpdate(updates, fmt.Sprintf(
		"optimizations set: FloatDecimals=%d ForceIPVersion=%v MaxPrefix=%d",
		dbConfig.Optimize.FloatDecimals,
		dbConfig.Optimize.ForceIPVersionEnabled(),
		dbConfig.Optimize.MaxPrefix,
	))
	sendUpdate(updates, fmt.Sprintf(
		"conditional resets: %+v",
		dbConfig.Merge.ConditionalResets,
	))

	// Close update channel when finished.
	if updates != nil {
		defer close(updates)
	}

	// Open output file to detect errors before processing.
	outputFile, err := os.Create(dbConfig.Output)
	if err != nil {
		return fmt.Errorf("failed to open output file for %s: %w", dbConfig.Name, err)
	}

	// Process sources.
	var (
		totalInserts   int
		totalStartTime = time.Now()
		slotStartTime  = time.Now()
	)
	for _, source := range sources {
		var inserted int
		sendUpdate(updates, fmt.Sprintf("---\nprocessing %s...", source.Name()))

		for {
			entry, err := source.NextEntry()
			if err != nil {
				sendUpdate(updates, fmt.Sprintf("failed to parse entry: %s", err.Error()))
				continue
			}
			if entry == nil {
				break
			}

			mmdbMap, err := entry.ToMMDBMap(dbConfig.Optimize)
			if err != nil {
				sendUpdate(updates, fmt.Sprintf("failed to convert %+v to mmdb map: %s", entry, err.Error()))
				continue
			}

			if entry.Net != nil {
				// Handle Network/Prefix Format.

				// Ignore entry if the IP version is forced and it does not match the mmdb DB.
				if dbConfig.Optimize.ForceIPVersionEnabled() && ipVersion(entry.Net.IP) != opts.IPVersion {
					continue
				}

				// Ignore entry if prefix is greater than the max prefix.
				if dbConfig.Optimize.MaxPrefix > 0 {
					prefixBits, _ := entry.Net.Mask.Size()
					if prefixBits > dbConfig.Optimize.MaxPrefix {
						continue
					}
				}

				err = writer.InsertFunc(entry.Net, ConditionalResetTopLevelMergeWith(mmdbMap, dbConfig.Merge.ConditionalResets))
				if err != nil {
					sendUpdate(updates, fmt.Sprintf("failed to insert %+v: %s", entry, err.Error()))
					continue
				}
			} else {
				// Handle From-To IP Format.

				// Ignore entry if the IP version is forced and it does not match the mmdb DB.
				if dbConfig.Optimize.ForceIPVersionEnabled() && ipVersion(entry.From) != opts.IPVersion {
					continue
				}

				start, ok1 := netip.AddrFromSlice(entry.From)
				end, ok2 := netip.AddrFromSlice(entry.To)
				if !ok1 || !ok2 {
					sendUpdate(updates, fmt.Sprintf("range with invalid IPs: %s - %s", entry.From, entry.To))
					continue
				}

				r := netipx.IPRangeFrom(start, end)
				if !r.IsValid() {
					sendUpdate(updates, fmt.Sprintf("range is invalid: %s - %s", entry.From, entry.To))
					continue
				}
				subnets := r.Prefixes()
				for _, subnet := range subnets {
					// Ignore entry if prefix is greater than the max prefix.
					if dbConfig.Optimize.MaxPrefix > 0 && subnet.Bits() > dbConfig.Optimize.MaxPrefix {
						continue
					}

					err = writer.InsertFunc(netipx.PrefixIPNet(subnet), ConditionalResetTopLevelMergeWith(mmdbMap, dbConfig.Merge.ConditionalResets))
					if err != nil {
						sendUpdate(updates, fmt.Sprintf("failed to insert %+v: %s", entry, err.Error()))
						continue
					}
				}
			}

			inserted++
			totalInserts++
			if inserted%reportSlotSize == 0 {
				sendUpdate(updates, fmt.Sprintf(
					"inserted %d entries - batch in %s (%s/op)",
					inserted,
					time.Since(slotStartTime).Round(time.Millisecond),
					(time.Since(slotStartTime)/reportSlotSize).Round(time.Microsecond),
				))
				slotStartTime = time.Now()
			}
		}
		if source.Err() != nil {
			return fmt.Errorf("source %s failed: %w", source.Name(), source.Err())
		}
		sendUpdate(updates, fmt.Sprintf(
			"inserted %d entries - batch in %s (%s/op)",
			inserted,
			time.Since(slotStartTime).Round(time.Millisecond),
			(time.Since(slotStartTime)/reportSlotSize).Round(time.Microsecond),
		))
	}

	// Write final db to file.
	_, err = writer.WriteTo(outputFile)
	if err != nil {
		return fmt.Errorf("faild to write %s to output file: %w", dbConfig.Name, err)
	}

	// Send final upate.
	var fileSize int64
	stat, err := os.Stat(dbConfig.Output)
	if err == nil {
		fileSize = stat.Size()
	}
	sendUpdate(updates, fmt.Sprintf(
		"---\n%s finished: inserted %d entries in %s, resulting in %.2f MB written to %s",
		dbConfig.Name,
		totalInserts,
		time.Since(totalStartTime).Round(time.Second),
		float64(fileSize)/1000000,
		dbConfig.Output,
	))

	return nil
}

// ConditionalResetTopLevelMergeWith is based on TopLevelMergeWith,
// but conditionally resets fields as defined in the conditional reset config.
// Both the new and existing value must be a Map. An error will be returned
// otherwise.
func ConditionalResetTopLevelMergeWith(newValue mmdbtype.DataType, cfg []ConditionalResetConfig) inserter.Func {
	return func(existingValue mmdbtype.DataType) (mmdbtype.DataType, error) {
		// Check if both values are maps before we start merging.
		newMap, ok := newValue.(mmdbtype.Map)
		if !ok {
			return nil, fmt.Errorf(
				"the new value is a %T, not a Map; ConditionalResetTopLevelMerge only works if both values are Map values",
				newValue,
			)
		}
		if existingValue == nil {
			return newValue, nil
		}
		existingMap, ok := existingValue.(mmdbtype.Map)
		if !ok {
			return nil, fmt.Errorf(
				"the existing value is a %T, not a Map; ConditionalResetTopLevelMerge only works if both values are Map values",
				existingValue,
			)
		}

		// Start merging.

		// First, do a normal top-level merge.
		returnMap := existingMap.Copy().(mmdbtype.Map) //nolint:forcetypeassert
		for k, v := range newMap {
			returnMap[k] = v.Copy()
		}

		// Then check which fields changed.
		for _, c := range cfg {
			var changed bool
			for _, key := range c.IfChanged {
				// Get existing value.
				existingSubVal, ok := existingMap[mmdbtype.String(key)]
				if !ok {
					// There is no existing value of that key, so there is no change possible.
					continue
				}
				// Get new value
				newSubVal, ok := newMap[mmdbtype.String(key)]
				if !ok {
					// Value of that key is not being set, so there is no change possible.
					continue
				}
				// Compare values if both are set.
				if !newSubVal.Equal(existingSubVal) {
					changed = true
					break
				}
			}
			// If any field changed, reset fields.
			if changed {
				for _, key := range c.Reset {
					resetVal, ok := newMap[mmdbtype.String(key)]
					if ok {
						// Reset with new value.
						returnMap[mmdbtype.String(key)] = resetVal
					} else {
						// Remove if no new value is present.
						delete(returnMap, mmdbtype.String(key))
					}
				}
			}
		}

		return returnMap, nil
	}
}

func sendUpdate(to chan string, msg string) {
	if to == nil {
		return
	}

	select {
	case to <- msg:
	default:
	}
}

func ipVersion(ip net.IP) int {
	if ip.To4() != nil {
		return 4
	}
	return 6
}
