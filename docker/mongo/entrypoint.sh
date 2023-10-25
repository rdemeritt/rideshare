#!/bin/bash

export MONGO_INITDB_ROOT_USERNAME="root"
export MONGO_INITDB_ROOT_PASSWORD="Password1!"

mongod --config /etc/mongod.conf