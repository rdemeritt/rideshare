package main

import (
	"googlemaps.github.io/maps"
)

func NewMapsClient() (*maps.Client, error) {
	client, err := maps.NewClient(maps.WithAPIKey(GoogleMapsAPIKey))
	check(err)

	return client, nil
}
