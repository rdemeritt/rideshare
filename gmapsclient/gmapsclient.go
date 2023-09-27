package gmapsclient

import (
	"rideshare/common"
	"googlemaps.github.io/maps"
)

func NewMapsClient() (*maps.Client, error) {
	client, err := maps.NewClient(maps.WithAPIKey(common.GoogleMapsAPIKey))
	common.Check(err)

	return client, nil
}
