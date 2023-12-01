package main

import "fmt"

// Location holds information regarding the geographical and network location of an IP address.
type Location struct {
	Continent struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"continent"`
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	Coordinates                  Coordinates `maxminddb:"location"`
	AutonomousSystemNumber       uint        `maxminddb:"autonomous_system_number"`
	AutonomousSystemOrganization string      `maxminddb:"autonomous_system_organization"`
	IsAnycast                    bool        `maxminddb:"is_anycast"`
	IsSatelliteProvider          bool        `maxminddb:"is_satellite_provider"`
	IsAnonymousProxy             bool        `maxminddb:"is_anonymous_proxy"`
}

// Coordinates holds geographic coordinates and their estimated accuracy.
type Coordinates struct {
	AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
	Latitude       float64 `maxminddb:"latitude"`
	Longitude      float64 `maxminddb:"longitude"`
}

type mmdbStats struct {
	Records             uint64
	Country             uint64
	Coordinates         uint64
	ASN                 uint64
	ASOrg               uint64
	IsAnycast           uint64
	IsSatelliteProvider uint64
	IsAnonymousProxy    uint64
}

func (ms *mmdbStats) Add(l *Location) {
	ms.Records++
	if l.Country.ISOCode != "" {
		ms.Country++
	}
	if l.Coordinates.Latitude != 0 && l.Coordinates.Longitude != 0 {
		ms.Coordinates++
	}
	if l.AutonomousSystemNumber != 0 {
		ms.ASN++
	}
	if l.AutonomousSystemOrganization != "" {
		ms.ASOrg++
	}
	if l.IsAnycast {
		ms.IsAnycast++
	}
	if l.IsSatelliteProvider {
		ms.IsSatelliteProvider++
	}
	if l.IsAnonymousProxy {
		ms.IsAnonymousProxy++
	}
}

func (ms *mmdbStats) Print() {
	fmt.Printf(
		"Total=%10d Country=%.2f%% Coords=%.2f%% ASN=%.2f%% ASOrg=%.2f%% AC=%.2f%% SP=%.2f%% AP=%.2f%%\n",
		ms.Records,
		(float64(ms.Country)/float64(ms.Records))*100,
		(float64(ms.Coordinates)/float64(ms.Records))*100,
		(float64(ms.ASN)/float64(ms.Records))*100,
		(float64(ms.ASOrg)/float64(ms.Records))*100,
		(float64(ms.IsAnycast)/float64(ms.Records))*100,
		(float64(ms.IsSatelliteProvider)/float64(ms.Records))*100,
		(float64(ms.IsAnonymousProxy)/float64(ms.Records))*100,
	)
}

type perDepthStats struct {
	layers [129]*mmdbStats
}

func (pds *perDepthStats) Add(depth int, l *Location) {
	ms := pds.layers[depth]
	if ms == nil {
		ms = &mmdbStats{}
		pds.layers[depth] = ms
	}
	ms.Add(l)
}

func (pds *perDepthStats) Print() {
	for depth, ms := range pds.layers {
		if ms == nil {
			continue
		}

		fmt.Printf("CIDR=/%3d ", depth)
		ms.Print()
	}
}
