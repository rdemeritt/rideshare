#!/bin/bash

eval $(cat /app/GMAPS_API_KEY | sed 's/^/export GMAPS_API_KEY=/')

# Load environment variables from files in /app/env
for file in /app/env/*; do
    name=$(basename "$file")
    value=$(cat "$file")
    export "$name"="$value"
done

if [ -n "$GMAPS_API_KEY" ]; then
    echo "GMAPS_API_KEY is set, removing secrets file"
    rm -f /app/GMAPS_API_KEY
else
    echo "GMAPS_API_KEY is not set"
fi

/app/trip_server -port 8080
