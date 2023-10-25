FROM mongo:latest

# Copy the MongoDB configuration file to the container
COPY mongod.conf /etc/mongod.conf

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose the MongoDB port
EXPOSE 27017

# Start MongoDB with the custom configuration file
# CMD ["mongod", "--config", "/etc/mongod.conf"]
ENTRYPOINT [ "/entrypoint.sh" ]
