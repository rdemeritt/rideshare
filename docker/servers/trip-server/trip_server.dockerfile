# Container image that runs your code
# FROM alpine:latest
# RUN apk add bash
FROM golang:1.21
RUN apt-get update

COPY trip_server/trip_server /app/trip_server
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /app/trip_server
RUN chmod +x /entrypoint.sh

RUN --mount=type=secret,id=GMAPS_API_KEY \
    cat /run/secrets/GMAPS_API_KEY > /app/GMAPS_API_KEY

RUN --mount=type=secret,id=RS_DB_USER \
    cat /run/secrets/RS_DB_USER > /app/RS_DB_USER

RUN --mount=type=secret,id=RS_DB_PASS \
    cat /run/secrets/RS_DB_PASS > /app/RS_DB_PASS

RUN --mount=type=secret,id=RS_DB_HOST \
    cat /run/secrets/RS_DB_HOST > /app/RS_DB_HOST

RUN --mount=type=secret,id=RS_DB_PORT \
    cat /run/secrets/RS_DB_PORT > /app/RS_DB_PORT

ENTRYPOINT [ "/entrypoint.sh" ]