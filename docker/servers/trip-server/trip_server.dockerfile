# Container image that runs your code
FROM alpine:latest
RUN apk add bash

COPY trip_server trip_server
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x trip_server
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]