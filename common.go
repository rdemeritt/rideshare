package main

const (
	// Replace with your own API key
	GoogleMapsAPIKey = "AIzaSyCBdtdaO3EjAgupwQCo0-IlOwxFW1w3UWk"

	// Google Maps API URL
	GoogleMapsAPIURL = "https://maps.googleapis.com/maps/api/directions/json?origin=%s&destination=%s&key=%s"
)

var (
	// Replace with your own source and destination coordinates
	SourceCoordinates      = "37.7749,-122.4194"
	DestinationCoordinates = "37.3352,-121.8811"
)
