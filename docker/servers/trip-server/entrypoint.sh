#!/bin/bash

eval $(cat /app/GMAPS_API_KEY | sed 's/^/export GMAPS_API_KEY=/')

if [ -n "$GMAPS_API_KEY" ]; then
    echo "GMAPS_API_KEY is set, removing secrets file"
    rm -f /app/GMAPS_API_KEY
else
    echo "GMAPS_API_KEY is not set"
fi

/app/trip_server -port 8080
