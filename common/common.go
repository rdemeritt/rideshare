package common

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// Replace with your own API key
	GoogleMapsAPIKey = "AIzaSyCBdtdaO3EjAgupwQCo0-IlOwxFW1w3UWk"
)

func Check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

func BsonMToJson(m bson.M) ([]byte, error) {
    // Convert the bson.M object to a map[string]interface{}
    var data map[string]interface{}
    bytes, err := bson.Marshal(m)
    if err != nil {
        return nil, err
    }
    err = bson.Unmarshal(bytes, &data)
    if err != nil {
        return nil, err
    }

    // Convert the map to JSON
    jsonBytes, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    return jsonBytes, nil
}

func BsonDToJson(d bson.D) ([]byte, error) {
    // Convert the bson.D object to a map[string]interface{}
    var data map[string]interface{}
    bytes, err := bson.Marshal(d)
    if err != nil {
        return nil, err
    }
    err = bson.Unmarshal(bytes, &data)
    if err != nil {
        return nil, err
    }

    // Convert the map to JSON
    jsonBytes, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    return jsonBytes, nil
}
