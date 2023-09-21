package main

import log "github.com/sirupsen/logrus"

const (
	// Replace with your own API key
	GoogleMapsAPIKey = "AIzaSyCBdtdaO3EjAgupwQCo0-IlOwxFW1w3UWk"

	// Google Maps API URL
	GoogleMapsAPIURL = "https://maps.googleapis.com/maps/api/directions/json?origin=%s&destination=%s&key=%s"
	routesGRPCUrl = "routes.googleapis.com:443"
	fieldMask  = "*" // fieldMask is a comma-separated list of fully qualified names of fields.
)

var (
	// Replace with your own source and destination coordinates
	SourceCoordinates      = "37.7749,-122.4194"
	DestinationCoordinates = "37.3352,-121.8811"
)

func check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}
