#!/bin/bash

# Load environment variables from files in /app/env
source ./env.sh

if [ -n "$GMAPS_API_KEY" ]; then
    echo "GMAPS_API_KEY is set"
else
    echo "GMAPS_API_KEY is not set"
fi

/app/trip_server -port 8080
