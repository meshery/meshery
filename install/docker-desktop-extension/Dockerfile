FROM golang:1.17-alpine AS builder
ENV CGO_ENABLED=0
RUN apk add --update make
WORKDIR /backend
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make bin

FROM node:14.17-alpine3.13 AS client-builder
WORKDIR /ui
# cache packages in layer
COPY ui-dde/package.json /ui/package.json
COPY ui-dde/package-lock.json /ui/package-lock.json
RUN --mount=type=cache,target=/usr/src/app/.npm \
    npm set cache /usr/src/app/.npm && \
    npm ci
# install
COPY ui-dde /ui
RUN npm run build

FROM alpine
LABEL org.opencontainers.image.title="meshery" \
    org.opencontainers.image.description="Meshery" \
    org.opencontainers.image.vendor="Layer5, Inc." \
    com.docker.desktop.extension.api.version=">= 0.2.0"
COPY --from=builder /backend/bin/service /
COPY docker-compose.yaml .
COPY metadata.json .
COPY meshery-logo-light.svg .
COPY --from=client-builder /ui/build ui
CMD /service -socket /run/guest-services/extension-meshery.sock
