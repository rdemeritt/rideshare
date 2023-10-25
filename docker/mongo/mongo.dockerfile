FROM mongo:latest

# Copy the MongoDB configuration file to the container
COPY mongod.conf /etc/mongod.conf

# Expose the MongoDB port
EXPOSE 27017

ENV MONGO_INITDB_ROOT_USERNAME="root"
ENV MONGO_INITDB_ROOT_PASSWORD="Password1!"

# Start MongoDB with the custom configuration file
CMD ["mongod", "--config", "/etc/mongod.conf"]
