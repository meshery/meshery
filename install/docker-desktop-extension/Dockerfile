FROM golang:1.16 as meshery-server
ARG TOKEN
ARG GIT_VERSION
ARG GIT_COMMITSHA
ARG RELEASE_CHANNEL

RUN adduser --disabled-login appuser
WORKDIR /github.com/meshery/meshery
ADD . .
RUN go clean -modcache; cd cmd; GOPROXY=https://proxy.golang.org GOSUMDB=off go build -ldflags="-w -s -X main.globalTokenForAnonymousResults=$TOKEN -X main.version=$GIT_VERSION -X main.commitsha=$GIT_COMMITSHA -X main.releasechannel=$RELEASE_CHANNEL" -tags draft -a -o /meshery .

FROM node:16 as ui
ADD ui-dde ui
RUN cd ui; npm install --only=production; npm run build && npm run export; mv out /

FROM node:16 as provider-ui
ADD provider-ui provider-ui
RUN cd provider-ui; npm install --only=production; npm run build && npm run export; mv out /

FROM frolvlad/alpine-glibc:alpine-3.13_glibc-2.32
LABEL org.opencontainers.image.title="Meshery" \
    org.opencontainers.image.description="Meshery extension for Docker Desktop" \
    org.opencontainers.image.vendor="Layer5" \
    com.docker.desktop.extension.api.version=">= 0.2.0"
    # com.docker.desktop.plugin.icon="meshery-logo-light.svg"
# NOTE (hershd23) : Commenting out the above line because of an error, will fix that
#RUN apt-get update; apt-get install -y ca-certificates; update-ca-certificates && rm -rf /var/lib/apt/lists/*
RUN apk update && apk add ca-certificates; update-ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates
COPY metadata.json .
COPY docker.svg .
COPY docker-compose.yaml .
COPY darwin /darwin
COPY ./oam /app/oam
COPY --from=ui /out /app/ui/out
