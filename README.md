# Build your own GeoIP .mmdb!

__Step 0: Setup Go__

Follow instructions here: <https://go.dev/dl/>

__Step 1: Compile__

    $ go build -C cmd/mmdbmeld
    $ go build -C cmd/mmdbcheck

__Step 2: Download geoip data sources__

_Default sources are CC0 and provided by <https://github.com/sapics/ip-location-db>_

    $ ./input/update.sh

    [...]

__Step 3: Build your MMDBs__

    $ ./cmd/mmdbmeld/mmdbmeld config-example.yml

    ==========
    building My IPv4 GeoIP DB
    database options set: IncludeReservedNetworks=true IPVersion=4 RecordSize=24
    optimizations set: FloatDecimals=2 ForceIPVersion=true MaxPrefix=0
    ---
    processing input/asn-ipv4.csv...
    inserted 100000 entries - batch in 314ms (3µs/op)
    inserted 200000 entries - batch in 287ms (3µs/op)
    inserted 300000 entries - batch in 281ms (3µs/op)
    inserted 368765 entries - batch in 206ms (2µs/op)
    ---
    processing input/geo-whois-asn-country-ipv4.csv...
    inserted 100000 entries - batch in 931ms (9µs/op)
    inserted 200000 entries - batch in 684ms (7µs/op)
    inserted 238714 entries - batch in 353ms (4µs/op)
    ---
    My IPv4 GeoIP DB finished: inserted 607479 entries in 3s, resulting in 7MB

    ==========
    building My IPv6 GeoIP DB
    database options set: IncludeReservedNetworks=true IPVersion=6 RecordSize=24
    optimizations set: FloatDecimals=2 ForceIPVersion=true MaxPrefix=0
    ---
    processing input/asn-ipv6.csv...
    inserted 86473 entries - batch in 382ms (4µs/op)
    ---
    processing input/geo-whois-asn-country-ipv6.csv...
    inserted 82381 entries - batch in 1.389s (14µs/op)
    ---
    My IPv6 GeoIP DB finished: inserted 168854 entries in 2s, resulting in 5MB

__Step 4: Check your MMDBs__

    $ ./cmd/mmdbcheck/mmdbcheck all output/geoip-v4.mmdb
    $ ./cmd/mmdbcheck/mmdbcheck all output/geoip-v6.mmdb
    
    [...]

    loading output/geoip-v4.mmdb with 7.85 MB

    Running all checks:

    Probing:
    ..............................................................................................................................................................................................................................................................::
    analyzed with 0 lookup errors
    Total=  14462461 Country=99.57% Coords=0.00% ASN=87.04% ASOrg=87.04% AC=0.00% SP=0.00% AP=0.00%

    Network Mask Stats:
              Total=    770040 Country=99.94% Coords=0.00% ASN=89.10% ASOrg=89.10% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/  7 Total=         1 Country=100.00% Coords=0.00% ASN=100.00% ASOrg=100.00% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/  8 Total=         6 Country=100.00% Coords=0.00% ASN=100.00% ASOrg=100.00% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/  9 Total=         9 Country=100.00% Coords=0.00% ASN=88.89% ASOrg=88.89% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 10 Total=        35 Country=100.00% Coords=0.00% ASN=88.57% ASOrg=88.57% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 11 Total=        90 Country=100.00% Coords=0.00% ASN=88.89% ASOrg=88.89% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 12 Total=       273 Country=100.00% Coords=0.00% ASN=89.74% ASOrg=89.74% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 13 Total=       561 Country=100.00% Coords=0.00% ASN=90.02% ASOrg=90.02% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 14 Total=      1321 Country=100.00% Coords=0.00% ASN=87.06% ASOrg=87.06% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 15 Total=      2743 Country=100.00% Coords=0.00% ASN=83.81% ASOrg=83.81% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 16 Total=     10414 Country=99.96% Coords=0.00% ASN=83.52% ASOrg=83.52% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 17 Total=      7469 Country=100.00% Coords=0.00% ASN=87.36% ASOrg=87.36% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 18 Total=     13200 Country=99.97% Coords=0.00% ASN=87.29% ASOrg=87.29% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 19 Total=     24993 Country=99.98% Coords=0.00% ASN=88.65% ASOrg=88.65% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 20 Total=     38550 Country=99.96% Coords=0.00% ASN=88.21% ASOrg=88.21% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 21 Total=     53388 Country=99.95% Coords=0.00% ASN=86.92% ASOrg=86.92% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 22 Total=    115769 Country=99.92% Coords=0.00% ASN=88.42% ASOrg=88.42% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 23 Total=    118221 Country=99.93% Coords=0.00% ASN=84.43% ASOrg=84.43% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 24 Total=    228493 Country=99.90% Coords=0.00% ASN=86.77% ASOrg=86.77% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 25 Total=     10503 Country=100.00% Coords=0.00% ASN=96.15% ASOrg=96.15% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 26 Total=     15323 Country=100.00% Coords=0.00% ASN=96.80% ASOrg=96.80% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 27 Total=     21747 Country=100.00% Coords=0.00% ASN=97.20% ASOrg=97.20% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 28 Total=     28184 Country=100.00% Coords=0.00% ASN=98.64% ASOrg=98.64% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 29 Total=     30933 Country=100.00% Coords=0.00% ASN=98.97% ASOrg=98.97% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 30 Total=     23054 Country=100.00% Coords=0.00% ASN=99.51% ASOrg=99.51% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 31 Total=      7488 Country=100.00% Coords=0.00% ASN=99.49% ASOrg=99.49% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 32 Total=     17272 Country=100.00% Coords=0.00% ASN=99.58% ASOrg=99.58% AC=0.00% SP=0.00% AP=0.00%

This check queries the following fields:

- Country: `country.iso_code`
- Coords: `location.latitude`, `location.longitude`
- ASN: `autonomous_system_number`
- ASOrg: `autonomous_system_organization`
- AC: `is_anycast`
- SP: `is_satellite_provider`
- AP: `is_anonymous_proxy`

### Customize your MMDBs.

Take a look at the <config-example.yml> to get an idea how to customize your MMDB.

There also several options to optimize your MMDB for smaller sizes (and reduced accuracy).

If you add more fields, it is a good idea to stick to the keys that MaxMind already uses. This keeps your MMDB compatible with existing systems.
You can find the keys they use here: https://pkg.go.dev/github.com/oschwald/geoip2-golang#City
