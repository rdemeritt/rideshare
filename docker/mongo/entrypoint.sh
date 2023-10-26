#!/bin/bash

eval $(export "MONGO_INITDB_ROOT_USERNAME=root")
eval $(export "MONGO_INITDB_ROOT_PASSWORD=Password1!")

mongod --config /etc/mongod.conf