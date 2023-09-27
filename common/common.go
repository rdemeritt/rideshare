package common

import log "github.com/sirupsen/logrus"

const (
	// Replace with your own API key
	GoogleMapsAPIKey = "***REMOVED***"
)

func Check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}
