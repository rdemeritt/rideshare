#!/bin/bash

eval $(cat /app/gmapsapikey | sed 's/^/export GMAPS_API_KEY=/')
echo $GMAPS_API_KEY
if [ -n "$GMAPS_API_KEY" ]; then
    echo "GMAPS_API_KEY is set, removing secrets file"
else
    echo "GMAPS_API_KEY is not set"
fi

/app/trip_server -port 8080
