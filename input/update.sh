#!/bin/bash

baseDir="$( cd "$(dirname "$0")" && pwd )"
cd "$baseDir"

echo "Downloading / updating geoip sources"
echo "Important: Always check the current license before using these sources!"
echo

##########

echo "updating iptoasn-asn-ipv4.csv"
curl --silent --show-error \
  --etag-save iptoasn-asn-ipv4.csv.etag \
  --etag-compare iptoasn-asn-ipv4.csv.etag \
  --output iptoasn-asn-ipv4.csv \
  "https://cdn.jsdelivr.net/npm/@ip-location-db/iptoasn-asn/iptoasn-asn-ipv4.csv"

echo "updating iptoasn-asn-ipv6.csv"
curl --silent --show-error \
  --etag-save iptoasn-asn-ipv6.csv.etag \
  --etag-compare iptoasn-asn-ipv6.csv.etag \
  --output iptoasn-asn-ipv6.csv \
  "https://cdn.jsdelivr.net/npm/@ip-location-db/iptoasn-asn/iptoasn-asn-ipv6.csv"

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
