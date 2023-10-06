package trip

import (
	"context"
	"fmt"

	"googlemaps.github.io/maps"
)

func PrintDirections(directions []maps.Route) {
	// Print the driving directions
	for _, route := range directions {
		for _, leg := range route.Legs {
			for _, step := range leg.Steps {
				fmt.Println(step.HTMLInstructions)
			}
		}
	}
}

func GetDirections(client *maps.Client, start string, end string) ([]maps.Route, error) {
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
