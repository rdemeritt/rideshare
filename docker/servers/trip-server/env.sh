#/bin/bash

# Load environment variables from files in /app/env
for file in /app/env/*; do
    name=$(basename "$file")
    value=$(cat "$file")
    eval $(echo export $name"="$value)
done
