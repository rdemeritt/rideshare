# Container image that runs your code
FROM alpine:latest
RUN apk add bash

COPY trip_server/trip_server /app/trip_server
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /app/trip_server
RUN chmod +x /entrypoint.sh

RUN --mount=type=secret,id=gmapsapikey \
    eval 'export GMAPS_API_KEY='$(cat /run/secrets/gmapsapikey)

ENTRYPOINT [ "/entrypoint.sh" ]