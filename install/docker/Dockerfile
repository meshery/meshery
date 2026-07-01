FROM golang:1.26.4 AS meshery-server
ARG TOKEN
ARG GIT_VERSION
ARG GIT_COMMITSHA
ARG RELEASE_CHANNEL

RUN adduser --disabled-login appuser
WORKDIR /github.com/meshery/meshery
ADD . .
RUN go clean -modcache; cd server; cd cmd; GOPROXY=https://proxy.golang.org GOSUMDB=off go build -ldflags="-w -s -X main.globalTokenForAnonymousResults=$TOKEN -X main.version=$GIT_VERSION -X main.commitsha=$GIT_COMMITSHA -X main.releasechannel=$RELEASE_CHANNEL" -tags draft -a -o /meshery .

# Compile the Go relationship/policy engine to WebAssembly and stage Go's
# wasm_exec.js runtime glue. ui/pages/_document.tsx loads
# /static/wasm/wasm_exec.js on every page, and the engine fetches
# /static/wasm/policy_engine.wasm at runtime. Both are .gitignore'd artifacts
# emitted by the `wasm-engine` Makefile target, which this container build never
# ran -- so the request fell through to the SPA HTML shell and the browser
# refused to execute it (MIME 'text/plain'). Full source is copied because
# server/policies/wasm/go.mod sets `replace github.com/meshery/meshery => ../../../`.
FROM golang:1.26.4 AS meshery-wasm
WORKDIR /github.com/meshery/meshery
ADD . .
RUN mkdir -p /wasm-dist && cd server/policies/wasm && \
    GOOS=js GOARCH=wasm GOPROXY=https://proxy.golang.org GOSUMDB=off \
        go build -trimpath -ldflags="-w -s" -o /wasm-dist/policy_engine.wasm . && \
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" /wasm-dist/wasm_exec.js && \
    awk '{ print; if (index($0, "chdir() { throw enosys(); },") > 0) { print "\t\t\tenv: {},"; found=1 } } END { exit(found ? 0 : 1) }' \
        /wasm-dist/wasm_exec.js > /wasm-dist/wasm_exec.js.tmp && \
    mv /wasm-dist/wasm_exec.js.tmp /wasm-dist/wasm_exec.js

FROM node:22-alpine AS base-node

# Installing NodeJS Dependencies
FROM base-node AS meshery-ui-deps
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY ui/package.json ui/package-lock.json ./
ENV HUSKY=0
RUN npm ci --legacy-peer-deps --ignore-scripts

FROM base-node AS meshery-ui-provider-deps
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY provider-ui/package.json provider-ui/package-lock.json ./
RUN npm ci --legacy-peer-deps

# Build NextJS Application
FROM base-node AS meshery-ui-builder
WORKDIR /app
COPY --from=meshery-ui-deps /app/node_modules ./node_modules
COPY ui .
# Stage the relationship/policy engine wasm + Go wasm runtime glue into the UI's
# public assets *before* the build so `next build` (output: 'export') emits them
# to out/static/wasm -- the path Meshery serves the UI from
# (server/router/server.go: ServeUI "../../ui/out/").
COPY --from=meshery-wasm /wasm-dist/ ./public/static/wasm/
# Disable NextJS Telemetry data collect
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

FROM base-node AS meshery-ui-provider-builder
WORKDIR /app
COPY --from=meshery-ui-provider-deps /app/node_modules ./node_modules
# Disable NextJS Telemetry data collect
ENV NEXT_TELEMETRY_DISABLED=1
COPY provider-ui .
RUN npm run build

# FROM alpine:3.23.2 AS seed_content
# RUN apk add --no-cache curl
# WORKDIR /
# Antiquated patterns; no longer compatible with the current version of Meshery
# RUN curl -L -s https://github.com/service-mesh-patterns/service-mesh-patterns/tarball/master -o service-mesh-patterns.tgz \
#     && mkdir service-mesh-patterns \
#     && mkdir -p /seed_content/patterns \
#     && tar xzf service-mesh-patterns.tgz --directory=service-mesh-patterns \
#     && mv service-mesh-patterns/*/samples/* /seed_content/patterns/

# Build-time Extension Point: This stage prepares optional .meshery directory from build context.
# System integrators can place a .meshery directory at the repository root to supply
# build-time configuration, data, packages, etc., which will be copied into the image.
FROM alpine:3.23.2 AS build_time_extensions
WORKDIR /build-ext
# Copy the entire context from meshery-server stage, then extract .meshery if it exists
COPY --from=meshery-server /github.com/meshery/meshery /meshery-context
RUN if [ -d /meshery-context/.meshery ]; then \
        cp -r /meshery-context/.meshery /build-ext/.meshery; \
    else \
        mkdir -p /build-ext/.meshery; \
    fi

# FROM alpine:3.23.2
# ENV GLIBC_REPO=https://github.com/sgerrand/alpine-pkg-glibc/releases/download
# ENV GLIBC_VERSION=2.35-r1
# RUN set -ex && \
#     apk --update add libstdc++ curl ca-certificates && \
#     apk add --no-cache --virtual .build-deps curl binutils zstd gcompat && \
#     for pkg in glibc-${GLIBC_VERSION} glibc-bin-${GLIBC_VERSION}; \
#     do curl -sSL ${GLIBC_REPO}/${GLIBC_VERSION}/${pkg}.apk -o /tmp/${pkg}.apk; done && \
#     apk del gcompat && \
#     apk add --force-overwrite --no-cache --allow-untrusted /tmp/*.apk && \
#     rm -v /tmp/*.apk && \
#     /usr/glibc-compat/sbin/ldconfig /lib /usr/glibc-compat/lib

FROM debian:bookworm-slim
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates

COPY ./server/cmd/server-config.env /app/server/cmd/server-config.env
COPY ./models /app/models
COPY ./server/permissions /app/server/permissions/
COPY ./docs/data/catalog /app/docs/data/catalog/
COPY --from=meshery-server /meshery /app/server/cmd/
COPY --from=meshery-server /etc/passwd /etc/passwd
COPY --from=meshery-ui-builder /app/out /app/ui/out
COPY --from=meshery-ui-builder /app/public /app/ui/public
COPY --from=meshery-ui-provider-builder /app/out /app/provider-ui/out
# COPY --from=seed_content /seed_content /home/appuser/.meshery/seed_content

# Permission keys are authorization keys which server will
# seed on initial boot for local provider, if there is update on this path
# ensure to update here also https://github.com/meshery/meshery/blob/85a0fd9b9cc03bbc51c4aa1efd1e837db72ef691/install/kubernetes/helm/meshery/values.yaml#L19
ENV KEYS_PATH="/app/server/permissions/keys.csv"

# Build-time Extension Point: Copy .meshery directory from build context if present.
# System integrators can place a .meshery directory at the repository root before running
# docker build to supply build-time configuration, data, packages, etc.
# The contents will be copied into /home/appuser/.meshery/ and merged with existing content.
COPY --from=build_time_extensions /build-ext/.meshery /home/appuser/.meshery

RUN addgroup docker
RUN mkdir -p /var/run; chown -R appuser:docker /var/run/
# Ensure .meshery and subdirectories exist with proper permissions
RUN mkdir -p /home/appuser/.meshery/config; chown -R appuser /home/appuser/; chown -R appuser /app/ui/public
USER appuser
WORKDIR /app/server/cmd
CMD ["./meshery"]
