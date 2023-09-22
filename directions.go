package main

import (
	"context"
	"fmt"

	"googlemaps.github.io/maps"
	log "github.com/sirupsen/logrus"
)

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

func getDistanceMatrix(client *maps.Client, driverLocation string, passengerStart string, passengerEnd string, units string) (*maps.DistanceMatrixResponse, error) {
    // Build the distance matrix request
    request := &maps.DistanceMatrixRequest{
        Origins:      []string{driverLocation, passengerStart},
        Destinations: []string{passengerStart, passengerEnd},
    }

	setUnits(units, request)

    // Send the distance matrix request
    response, err := client.DistanceMatrix(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response, nil
}

func setUnits(units string, r *maps.DistanceMatrixRequest) {
	switch units {
	case "metric":
		r.Units = maps.UnitsMetric
	case "imperial":
		r.Units = maps.UnitsImperial
	case "":
		// ignore
	default:
		log.Fatalf("Unknown units %s", units)
	}
}

func printFullDistanceMatrix(resp *maps.DistanceMatrixResponse) {
	// Print the distance matrix
    for _, row := range resp.Rows {
		fmt.Println("row: ", row)
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

func printDistanceMatrix(resp *maps.DistanceMatrixResponse) {
    // Print the distance and duration for the first and last elements
    if len(resp.Rows) > 0 && len(resp.Rows[0].Elements) > 0 {
        // Print the distance and duration for the first element
        element := resp.Rows[0].Elements[0]
        if element.Status == "OK" {
            log.Debugf("Distance for first element: %v\n", element.Distance.HumanReadable)
            log.Debugf("Duration for first element: %v\n", element.Duration)
        } else {
            log.Debugf("Error for first element: %v\n", element.Status)
        }

        // Print the distance and duration for the last element
        element = resp.Rows[len(resp.Rows)-1].Elements[len(resp.Rows[0].Elements)-1]
        if element.Status == "OK" {
            log.Debugf("Distance for last element: %v\n", element.Distance.HumanReadable)
            log.Debugf("Duration for last element: %v\n", element.Duration)
        } else {
            log.Debugf("Error for last element: %v\n", element.Status)
        }
    }
}