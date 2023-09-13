package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type DirectionsResponse struct {
	Routes []struct {
		Legs []struct {
			Steps []struct {
				HtmlInstructions string `json:"html_instructions"`
			} `json:"steps"`
		} `json:"legs"`
	} `json:"routes"`
}

func setStart(start string) {
	SourceCoordinates = start
}

func setEnd(end string) {
	DestinationCoordinates = end
}

func printDirections(directions DirectionsResponse) {
	// Print the driving directions
	for _, route := range directions.Routes {
		for _, leg := range route.Legs {
			for _, step := range leg.Steps {
				log.Info(step.HtmlInstructions)
			}
		}
	}
}

func getDirections(begin string, end string) DirectionsResponse {
	// Build the API request URL
	url := fmt.Sprintf(GoogleMapsAPIURL, begin, end, GoogleMapsAPIKey)
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
