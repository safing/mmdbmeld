# Build your own GeoIP .mmdb!

__Step 0: Setup Go__

Follow instructions here: <https://go.dev/dl/>

__Step 1: Compile__

    $ go build -C cmd/mmdbmeld
    $ go build -C cmd/mmdbcheck
    $ go build -C cmd/mmdbquery

__Step 2: Download geoip data sources__

_Default sources are CC0/PDDL and provided by <https://github.com/sapics/ip-location-db>_

    $ ./input/update.sh

    [...]

Note:  
If you start using mmdbmeld, it is highly recommended you choose geoip data sources suitable for your use case, in terms of size and quality.
Be sure to also check the licenses of the data sources before using them.

__Step 3: Build your MMDBs__

    $ ./cmd/mmdbmeld/mmdbmeld config-example.yml

    ==========
    building My IPv4 GeoIP DB
    database options set: IPVersion=4 RecordSize=24 (IncludeReservedNetworks=true DisableIPv4Aliasing=true)
    database types: autonomous_system_number, autonomous_system_organization, country.iso_code
    optimizations set: FloatDecimals=2 ForceIPVersion=true MaxPrefix=0
    conditional resets: [{IfChanged:[country] Reset:[location]}]
    ---
    processing input/iptoasn-asn-ipv4.csv...
    inserted 100000 entries - batch in 356ms (4µs/op)
    inserted 200000 entries - batch in 327ms (3µs/op)
    inserted 300000 entries - batch in 322ms (3µs/op)
    inserted 360970 entries - batch in 224ms (2µs/op)
    ---
    processing input/geo-whois-asn-country-ipv4.csv...
    inserted 100000 entries - batch in 1.012s (10µs/op)
    inserted 200000 entries - batch in 889ms (9µs/op)
    inserted 238756 entries - batch in 418ms (4µs/op)
    ---
    My IPv4 GeoIP DB finished: inserted 599726 entries in 4s, resulting in 7.67 MB written to output/geoip-v4.mmdb

    ==========
    building My IPv6 GeoIP DB
    database options set: IPVersion=6 RecordSize=24 (IncludeReservedNetworks=true DisableIPv4Aliasing=true)
    database types: autonomous_system_number, autonomous_system_organization, country.iso_code
    optimizations set: FloatDecimals=2 ForceIPVersion=true MaxPrefix=0
    conditional resets: [{IfChanged:[country] Reset:[location]}]
    ---
    processing input/iptoasn-asn-ipv6.csv...
    inserted 81939 entries - batch in 352ms (4µs/op)
    ---
    processing input/geo-whois-asn-country-ipv6.csv...
    inserted 82321 entries - batch in 1.609s (16µs/op)
    ---
    My IPv6 GeoIP DB finished: inserted 164260 entries in 2s, resulting in 5.42 MB written to output/geoip-v6.mmdb

__Step 4: Check your MMDBs__

    $ ./cmd/mmdbcheck/mmdbcheck all output/geoip-v4.mmdb
    $ ./cmd/mmdbcheck/mmdbcheck all output/geoip-v6.mmdb
    
    [...]

    loading output/geoip-v4.mmdb with 7.67 MB

    Running all checks:

    Probing:
    ..............................................................................................................................................................................................................................................................::
    analyzed with 0 lookup errors
    Total=  14462461 Country=99.57% Coords=0.00% ASN=82.40% ASOrg=82.40% AC=0.00% SP=0.00% AP=0.00%

    Network Mask Stats:
            Total=    800822 Country=99.93% Coords=0.00% ASN=85.24% ASOrg=85.24% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/  7 Total=         1 Country=100.00% Coords=0.00% ASN=100.00% ASOrg=100.00% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/  8 Total=         6 Country=100.00% Coords=0.00% ASN=100.00% ASOrg=100.00% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/  9 Total=         8 Country=100.00% Coords=0.00% ASN=62.50% ASOrg=62.50% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 10 Total=        35 Country=100.00% Coords=0.00% ASN=71.43% ASOrg=71.43% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 11 Total=        85 Country=100.00% Coords=0.00% ASN=78.82% ASOrg=78.82% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 12 Total=       278 Country=100.00% Coords=0.00% ASN=83.45% ASOrg=83.45% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 13 Total=       569 Country=100.00% Coords=0.00% ASN=87.17% ASOrg=87.17% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 14 Total=      1327 Country=100.00% Coords=0.00% ASN=84.78% ASOrg=84.78% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 15 Total=      2760 Country=100.00% Coords=0.00% ASN=82.14% ASOrg=82.14% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 16 Total=     10167 Country=99.95% Coords=0.00% ASN=80.14% ASOrg=80.14% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 17 Total=      7563 Country=100.00% Coords=0.00% ASN=83.49% ASOrg=83.49% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 18 Total=     13354 Country=99.97% Coords=0.00% ASN=83.91% ASOrg=83.91% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 19 Total=     25354 Country=99.98% Coords=0.00% ASN=84.98% ASOrg=84.98% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 20 Total=     39116 Country=99.96% Coords=0.00% ASN=84.53% ASOrg=84.53% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 21 Total=     54874 Country=99.95% Coords=0.00% ASN=82.60% ASOrg=82.60% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 22 Total=    117368 Country=99.94% Coords=0.00% ASN=83.46% ASOrg=83.46% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 23 Total=    124541 Country=99.94% Coords=0.00% ASN=79.46% ASOrg=79.46% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 24 Total=    240113 Country=99.92% Coords=0.00% ASN=81.46% ASOrg=81.46% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 25 Total=     10581 Country=99.94% Coords=0.00% ASN=95.53% ASOrg=95.53% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 26 Total=     15500 Country=99.95% Coords=0.00% ASN=96.35% ASOrg=96.35% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 27 Total=     22059 Country=99.91% Coords=0.00% ASN=96.78% ASOrg=96.78% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 28 Total=     30393 Country=99.91% Coords=0.00% ASN=98.38% ASOrg=98.38% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 29 Total=     34128 Country=99.89% Coords=0.00% ASN=98.76% ASOrg=98.76% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 30 Total=     24836 Country=99.81% Coords=0.00% ASN=99.44% ASOrg=99.44% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 31 Total=      7847 Country=99.86% Coords=0.00% ASN=98.89% ASOrg=98.89% AC=0.00% SP=0.00% AP=0.00%
    CIDR=/ 32 Total=     17959 Country=99.96% Coords=0.00% ASN=99.20% ASOrg=99.20% AC=0.00% SP=0.00% AP=0.00%

This check queries the following fields:

- Country: `country.iso_code`
- Coords: `location.latitude`, `location.longitude`
- ASN: `autonomous_system_number`
- ASOrg: `autonomous_system_organization`
- AC: `is_anycast`
- SP: `is_satellite_provider`
- AP: `is_anonymous_proxy`

__Step 5: Query your MMDBs__

    $ ./cmd/mmdbquery/mmdbquery output/geoip-v4.mmdb 1.1.1.1

    1.1.1.0/24:
      autonomous_system_number: 13335
      autonomous_system_organization: CLOUDFLARENET
      country:
        iso_code: AU

### Customize your MMDBs.

Take a look at the [config-example.yml](https://github.com/safing/mmdbmeld/blob/master/config-example.yml) to get an idea how to customize your MMDB.

There also several options to optimize your MMDB for smaller sizes (and reduced accuracy).

It is highly recommended you choose geoip data sources suitable for your use case, in terms of size and quality.
Be sure to also check the licenses of the data sources before using them.

If you add more fields, it is a good idea to stick to the keys that MaxMind already uses. This keeps your MMDB compatible with existing systems.
You can find the keys they use here: https://pkg.go.dev/github.com/oschwald/geoip2-golang#City

### Supported Sources


Start by defining the fields and their data type of the output mmdb:

```yaml
databases:
  - name: "My IPv4 GeoIP DB"
    types:
      "country.iso_code": string
      "location.latitude": float32
      "location.longitude": float32
      "autonomous_system_organization": string
      "autonomous_system_number": uint32
      "is_anycast": bool
      "is_satellite_provider": bool
      "is_anonymous_proxy": bool
```

There are three special fields which are not defined in the types:

- from: The start address of an IP range.
- to: The end address of an IP range.
- net: An IP range in CIDR notation.

These are used to derive the IP ranges the data (row, entry) is applicable for.

##### CSV

Define columns with `fields`, which must match a field defined in the `types`:

```yaml
databases:
  ...
    fields: ["from", "to", "autonomous_system_number", "autonomous_system_organization"]
```

All rows must have exactly the specified amount of columns. Use `-` to define a column you are not using, eg.:

```yaml
databases:
  ...
    fields: ["from", "to", "country.iso_code", "-", "-", "-", "-", "location.latitude", "location.longitude", "-"]
```

##### IPFire

The [IPFire Firewall](https://www.ipfire.org/) maintains a [geoip database in a custom format](https://git.ipfire.org/?p=location/location-database.git;a=summary), which notably includes IP categorization, such as `is-anycast`.

Define fields with `fieldMap`, mapping IPFire database keys to `types`:

```yaml
databases:
  ...
    fieldMap:
      "aut-num": "autonomous_system_number"
      "name": "autonomous_system_organization"
      "country": "country.iso_code"
      "is-anycast": "is_anycast"
      "is-satellite-provider": "is_satellite_provider"
      "is-anonymous-proxy": "is_anonymous_proxy"
```

### Defaults

If you are building more than one mmdb file, you can use defaults to apply certain configuration to all databases.
A common case is when you build an IPv4 and IPv6 database separately.

When using defaults, always check if they are applied correctly by checking the logs when building the databases.

```yaml
defaults:
  types: # Entries are merged.
    # Databases inherit defaults types, if not yet set.
    # You can define a type with a "-" type string to ignore a default in a database config.
    "country.iso_code": string
    "location.latitude": float32
    "location.longitude": float32
    "autonomous_system_organization": string
    "autonomous_system_number": uint32
    "is_anycast": bool
    "is_satellite_provider": bool
    "is_anonymous_proxy": bool
  optimize: # Entries are used as default separately.
    floatDecimals: 2 # Default is used when database value is 0.
    forceIPVersion: true # Default is used when database value is not defined.
    maxPrefix: 24 # Default is used when database value is 0.
  merge: # Entries are used as default separately.
    conditionalResets: # Default is used when not defined or empty in database config.
    - ifChanged: ["country"]
        reset: ["location"]
```
