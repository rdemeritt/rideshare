package gmapsclient

import (
	"googlemaps.github.io/maps"
	"rideshare/common"
)

func NewMapsClient() (*maps.Client, error) {
	client, err := maps.NewClient(maps.WithAPIKey(common.GoogleMapsAPIKey))
	common.Check(err)

	return client, nil
}
