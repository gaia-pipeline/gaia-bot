FROM alpine
ARG GOARCH=amd64

RUN apk add -u ca-certificates
ADD ./bin/linux/${GOARCH}/gaia-bot /app/
ADD ./scripts/push_tag.sh /usr/local/push_tag.sh

WORKDIR /app/
ENTRYPOINT [ "/app/gaia-bot" ]
