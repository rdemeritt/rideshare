package gmapsclient

import (
	"rideshare/common"

	"googlemaps.github.io/maps"
)

func NewMapsClient(key string) (*maps.Client, error) {
	client, err := maps.NewClient(maps.WithAPIKey(key))
	common.Check(err)

	return client, nil
}
