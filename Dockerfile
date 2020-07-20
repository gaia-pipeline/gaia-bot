FROM golang:1.13.8-stretch
ARG GOARCH=amd64

RUN apt-get update && apt-get install -y git build-essential apt-transport-https ca-certificates curl gnupg-agent software-properties-common
RUN printf '#!/bin/sh\nexit 0' > /usr/sbin/policy-rc.d
RUN /etc/init.d/dbus start
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
RUN apt-key fingerprint 0EBFCD88
RUN add-apt-repository \
       "deb [arch=amd64] https://download.docker.com/linux/debian \
       $(lsb_release -cs) \
       stable"
RUN apt-get update
RUN apt-get install -y docker-ce docker-ce-cli

ADD ./build/linux/${GOARCH}/gaia-bot /app/
ADD ./scripts/fetch_pr.sh /usr/local/bin
ADD ./scripts/push_tag.sh /usr/local/bin

EXPOSE 9998

WORKDIR /app/
ENTRYPOINT [ "/app/gaia-bot" ]
