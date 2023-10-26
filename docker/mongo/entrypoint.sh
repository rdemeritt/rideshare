#!/bin/bash

echo "root" > /app/env/MONGO_INITDB_ROOT_USERNAME
echo "Password1!" > /app/env/MONGO_INITDB_ROOT_PASSWORD

source /app/env.sh

mongod --config /etc/mongod.conf