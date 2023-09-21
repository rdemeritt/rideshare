package main

import (
	"context"
	"fmt"

	"googlemaps.github.io/maps"
	log "github.com/sirupsen/logrus"
)

func setStart(start string) {
	SourceCoordinates = start
}

func setEnd(end string) {
	DestinationCoordinates = end
}

func printDirections(directions []maps.Route) {
	// Print the driving directions
    for _, route := range directions {
        for _, leg := range route.Legs {
            for _, step := range leg.Steps {
                fmt.Println(step.HTMLInstructions)
            }
        }
    }
}

func getDirections(client maps.Client, begin string, end string) []maps.Route {
	// Build the API request URL
	url := fmt.Sprintf(GoogleMapsAPIURL, begin, end, GoogleMapsAPIKey)
	log.Debugf("API request URL: %s", url)

	 // Build the directions request
	 	directionsRequest := &maps.DirectionsRequest{
        Origin:      begin,
        Destination: end,
    }

	// Send the directions request
    routes, _, err := client.Directions(context.Background(), directionsRequest)
    check(err)

	return routes
}
