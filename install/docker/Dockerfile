FROM golang:1.25.5 AS meshery-server
ARG TOKEN
ARG GIT_VERSION
ARG GIT_COMMITSHA
ARG RELEASE_CHANNEL

WORKDIR /src

# Cache Go modules (best practice for faster rebuilds)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -ldflags="-w -s \
      -X main.globalTokenForAnonymousResults=$TOKEN \
      -X main.version=$GIT_VERSION \
      -X main.commitsha=$GIT_COMMITSHA \
      -X main.releasechannel=$RELEASE_CHANNEL" \
    -tags draft -a \
    -o /meshery ./server/cmd

FROM node:22-alpine AS base-node

# Installing NodeJS Dependencies
FROM base-node AS meshery-ui-deps
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY ui/package.json ui/package-lock.json ./
ENV HUSKY=0
RUN npm ci --omit=dev --legacy-peer-deps --ignore-scripts

FROM base-node AS meshery-ui-provider-deps
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY provider-ui/package.json provider-ui/package-lock.json ./
RUN npm ci --omit=dev --legacy-peer-deps --ignore-scripts

# Build NextJS Application
FROM base-node AS meshery-ui-builder
WORKDIR /app
COPY --from=meshery-ui-deps /app/node_modules ./node_modules
COPY ui .
# Disable NextJS Telemetry data collection
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

FROM base-node AS meshery-ui-provider-builder
WORKDIR /app
COPY --from=meshery-ui-provider-deps /app/node_modules ./node_modules
# Disable NextJS Telemetry data collection
ENV NEXT_TELEMETRY_DISABLED=1
COPY provider-ui .
RUN npm run build

FROM alpine:3.23.4 AS seed_content
RUN apk add --no-cache curl
WORKDIR /
# Download latest WASM filters (robust, fixed syntax)
RUN mkdir -p seed_content/filters/binaries && \
    curl -s https://api.github.com/repos/layer5io/wasm-filters/releases/latest | \
    grep '"browser_download_url":.*\.wasm"' | \
    cut -d '"' -f 4 | \
    xargs -I {} curl -L -o seed_content/filters/binaries/$(basename {}) {}

# Build-time Extension Point: This stage prepares optional .meshery directory from build context.
# System integrators can place a .meshery directory at the repository root to supply
# build-time configuration, data, packages, etc., which will be copied into the image.
FROM alpine:3.23.4 AS build_time_extensions
WORKDIR /build-ext
# Copy the entire context from meshery-server stage, then extract .meshery if it exists
COPY --from=meshery-server /src /meshery-context
RUN if [ -d /meshery-context/.meshery ]; then \
        cp -r /meshery-context/.meshery /build-ext/.meshery; \
    else \
        mkdir -p /build-ext/.meshery; \
    fi

FROM debian:bookworm-slim
ENV DEBIAN_FRONTEND=noninteractive

# === Combined RUN for best practices (single layer, proper permissions, user creation) ===
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/* && \
    update-ca-certificates && \
    # Create non-root user + docker group (for socket access)
    groupadd -r docker && \
    useradd -r -u 1000 -g docker -s /usr/sbin/nologin -c "Meshery runtime user" appuser && \
    # Create required directories
    mkdir -p /var/run /home/appuser/.meshery/config /app/ui/public && \
    # Fix ownership (files copied earlier are owned by root)
    chown -R appuser:docker /var/run/ && \
    chown -R appuser:appuser /home/appuser/ /app/ui/public

# Copy application artifacts
COPY ./server/cmd/server-config.env /app/server/cmd/server-config.env
COPY ./server/meshmodel /app/server/meshmodel
COPY ./server/permissions /app/server/permissions/
COPY ./docs/data/catalog /app/docs/data/catalog/

COPY --from=meshery-server /meshery /app/server/cmd/

COPY --from=meshery-ui-builder /app/out /app/ui/out
COPY --from=meshery-ui-builder /app/public /app/ui/public
COPY --from=meshery-ui-provider-builder /app/out /app/provider-ui/out

# Build-time Extension Point: Copy .meshery directory from build context if present first
# (optional overrides/config), then overlay the always-present seed content
COPY --from=build_time_extensions /build-ext/.meshery /home/appuser/.meshery
COPY --from=seed_content /seed_content /home/appuser/.meshery/seed_content

# Permission keys are authorization keys which server will
# seed on initial boot for local provider
ENV KEYS_PATH="/app/server/permissions/keys.csv"

# OCI labels (best practice)
LABEL org.opencontainers.image.source="https://github.com/meshery/meshery" \
      org.opencontainers.image.title="Meshery" \
      org.opencontainers.image.description="Meshery - the cloud native manager" \
      org.opencontainers.image.version="${GIT_VERSION}" \
      org.opencontainers.image.revision="${GIT_COMMITSHA}"

# Meshery server default port
EXPOSE 8080

# Health check (best practice)
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

USER appuser
WORKDIR /app/server/cmd
CMD ["./meshery"]