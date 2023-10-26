FROM mongo:latest

# Copy the MongoDB configuration file to the container
COPY mongod.conf /etc/mongod.conf

COPY env.sh /app/env.sh
RUN chmod +x /app/env.sh
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Expose the MongoDB port
EXPOSE 27017

# Start MongoDB with the custom configuration file
ENTRYPOINT [ "/app/entrypoint.sh" ]
