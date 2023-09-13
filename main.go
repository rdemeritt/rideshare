package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func init() {
	initLog()
	parseCommandLineFlags()
}

func parseCommandLineFlags() {
	// Define command line flags
	logLevel := flag.String("log-level", "info", "Logging level (debug, info, warn, error)")
	startCoord := flag.String("start", "37.7749,-122.4194", "Starting coordinates")
	endCoord := flag.String("end", "37.3352,-121.8811", "Ending coordinates")

	// Parse command line flags
	flag.Parse()

	// Set the logging level
	setLogLevel(*logLevel)
	// Set the start and end coordinates
	setStart(*startCoord)
	setEnd(*endCoord)
}

func setStart(start string) {
	SourceCoordinates = start
}

func setEnd(end string) {
	DestinationCoordinates = end
}

func getDirections(begin string, end string) DirectionsResponse {
    // Build the API request URL
	url := fmt.Sprintf(GoogleMapsAPIURL, SourceCoordinates, DestinationCoordinates, GoogleMapsAPIKey)
	log.Debugf("API request URL: %s", url)

    // Send the API request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read API response body: %v", err)
	}
	log.Debug(string(body))

	if resp.StatusCode != http.StatusOK {

		log.Fatalf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var directions DirectionsResponse
	err = json.Unmarshal(body, &directions)
	if err != nil {
		log.Fatalf("Failed to parse JSON response: %v", err)
	}

    return directions
}

func main() {
    // Get the driving directions
    directions := getDirections(SourceCoordinates, DestinationCoordinates)

	// Print the driving directions
	for _, route := range directions.Routes {
		for _, leg := range route.Legs {
			for _, step := range leg.Steps {
				log.Info(step.HtmlInstructions)
			}
		}
	}
}
