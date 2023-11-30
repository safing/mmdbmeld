#!/bin/bash

baseDir="$( cd "$(dirname "$0")" && pwd )"
cd "$baseDir"

echo "Downloading / updating geoip sources"
echo "Important: Always check the current license before using these sources!"
echo

##########

echo "updating asn-ipv4.csv"
curl --silent --show-error \
  --etag-save asn-ipv4.csv.etag \
  --etag-compare asn-ipv4.csv.etag \
  --output asn-ipv4.csv \
  "https://cdn.jsdelivr.net/npm/@ip-location-db/asn/asn-ipv4.csv"

echo "updating asn-ipv6.csv"
curl --silent --show-error \
  --etag-save asn-ipv6.csv.etag \
  --etag-compare asn-ipv6.csv.etag \
  --output asn-ipv6.csv \
  "https://cdn.jsdelivr.net/npm/@ip-location-db/asn/asn-ipv6.csv"

##########

echo "updating geo-whois-asn-country-ipv4.csv"
curl --silent --show-error \
  --etag-save geo-whois-asn-country-ipv4.csv.etag \
  --etag-compare geo-whois-asn-country-ipv4.csv.etag \
  --output geo-whois-asn-country-ipv4.csv \
  "https://cdn.jsdelivr.net/npm/@ip-location-db/geo-whois-asn-country/geo-whois-asn-country-ipv4.csv"

echo "updating geo-whois-asn-country-ipv6.csv"
curl --silent --show-error \
  --etag-save geo-whois-asn-country-ipv6.csv.etag \
  --etag-compare geo-whois-asn-country-ipv6.csv.etag \
  --output geo-whois-asn-country-ipv6.csv \
  "https://cdn.jsdelivr.net/npm/@ip-location-db/geo-whois-asn-country/geo-whois-asn-country-ipv6.csv"
