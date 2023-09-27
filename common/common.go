package common

import log "github.com/sirupsen/logrus"

const (
	// Replace with your own API key
	GoogleMapsAPIKey = "AIzaSyCBdtdaO3EjAgupwQCo0-IlOwxFW1w3UWk"
)

func Check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}
