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

func getDirections(client *maps.Client, start string, end string) ([]maps.Route, error) {
	 // Build the directions request
	 	directionsRequest := &maps.DirectionsRequest{
        Origin:      start,
        Destination: end,
    }

	// Send the directions request
    routes, _, err := client.Directions(context.Background(), directionsRequest)
    if err != nil {
        return nil, err
    }

	return routes, nil
}

func getDistanceMatrix(client *maps.Client, start string, end string) (*maps.DistanceMatrixResponse, error) {
    // Build the distance matrix request
    req := &maps.DistanceMatrixRequest{
        Origins:      []string{start},
        Destinations: []string{end},
    }

    // Send the distance matrix request
    resp, err := client.DistanceMatrix(context.Background(), req)
    if err != nil {
        return nil, err
    }

    return resp, nil
}

func printDistanceMatrix(resp *maps.DistanceMatrixResponse) {
	// Print the distance matrix
    for _, row := range resp.Rows {
        for _, element := range row.Elements {
            if element.Status == "OK" {
                log.Debugf("Distance: %v\n", element.Distance.HumanReadable)
                log.Debugf("Duration: %v\n", element.Duration)
            } else {
                log.Debugf("Error: %v\n", element.Status)
            }
        }
    }
}