FROM alpine
ARG GOARCH=amd64

RUN apk add -u ca-certificates
ADD ./bin/linux/${GOARCH}/gaia-bot /app/

WORKDIR /app/
ENTRYPOINT [ "/app/gaia-bot" ]
