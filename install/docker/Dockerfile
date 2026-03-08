FROM golang:1.25.5 AS meshery-server
ARG TOKEN
ARG GIT_VERSION
ARG GIT_COMMITSHA
ARG RELEASE_CHANNEL

WORKDIR /github.com/meshery/meshery
ADD . .
RUN make -f install/docker/Makefile build-meshery-server TOKEN="$TOKEN" GIT_VERSION="$GIT_VERSION" GIT_COMMITSHA="$GIT_COMMITSHA" RELEASE_CHANNEL="$RELEASE_CHANNEL"

FROM node:20 AS base-node

# Installing NodeJS Dependencies
FROM base-node AS meshery-ui-deps
WORKDIR /app
COPY ui/package.json ui/package-lock.json ./
COPY install/docker/Makefile /app/Makefile
ENV HUSKY=0
RUN make -f /app/Makefile ui-deps

FROM base-node AS meshery-ui-provider-deps
WORKDIR /app
COPY provider-ui/package.json provider-ui/package-lock.json ./
COPY install/docker/Makefile /app/Makefile
RUN make -f /app/Makefile provider-ui-deps

# Build NextJS Application
FROM base-node AS meshery-ui-builder
WORKDIR /app
COPY --from=meshery-ui-deps /app/node_modules ./node_modules
COPY install/docker/Makefile /app/Makefile
COPY ui .
# Disable NextJS Telemetry data collect
ENV NEXT_TELEMETRY_DISABLED=1
RUN make -f /app/Makefile ui-build

FROM base-node AS meshery-ui-provider-builder
WORKDIR /app
COPY --from=meshery-ui-provider-deps /app/node_modules ./node_modules
COPY install/docker/Makefile /app/Makefile
# Disable NextJS Telemetry data collect
ENV NEXT_TELEMETRY_DISABLED=1
COPY provider-ui .
RUN make -f /app/Makefile provider-ui-build

# FROM ubuntu AS wrk2
# RUN apt-get -y update && apt-get -y install git && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/man/?? /usr/share/man/??_*
# RUN apt-get -y update && apt-get -y  install build-essential libssl-dev git zlib1g-dev
# RUN git config --global user.email "meshery@meshery.io"
# RUN git config --global user.name "meshery"
# RUN git clone --depth=1 https://github.com/layer5io/wrk2 && cd wrk2 && make

FROM debian:bookworm-slim AS seed_content
COPY --from=base-node /usr/bin/make /usr/bin/make
COPY install/docker/Makefile /Makefile
WORKDIR /
RUN make -f /Makefile install-seed-content-deps
RUN make -f /Makefile seed-filters
# Antiquated patterns; no longer compatible with the current version of Meshery
# RUN curl -L -s https://github.com/service-mesh-patterns/service-mesh-patterns/tarball/master -o service-mesh-patterns.tgz \
#     && mkdir service-mesh-patterns \
#     && mkdir -p /seed_content/patterns \
#     && tar xzf service-mesh-patterns.tgz --directory=service-mesh-patterns \
#     && mv service-mesh-patterns/*/samples/* /seed_content/patterns/

#FROM ubuntu AS nighthawk
#RUN apt-get -y update && apt-get -y install git && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/man/?? /usr/share/man/??_*
#RUN apt-get -y update && apt-get -y  install build-essential libssl-dev git zlib1g-dev
#RUN git config --global user.email "ci@meshery.io"
#RUN git config --global user.name "meshery"
#RUN git clone https://github.com/layer5io/nighthawk-go
#RUN cd nighthawk-go/apinighthawk/bin && chmod +x ./nighthawk_client

# Build-time Extension Point: This stage prepares optional .meshery directory from build context.
# System integrators can place a .meshery directory at the repository root to supply
# build-time configuration, data, packages, etc., which will be copied into the image.
FROM debian:bookworm-slim AS build_time_extensions
COPY --from=base-node /usr/bin/make /usr/bin/make
COPY install/docker/Makefile /Makefile
WORKDIR /build-ext
# Copy the entire context from meshery-server stage, then extract .meshery if it exists
COPY --from=meshery-server /github.com/meshery/meshery /meshery-context
RUN make -f /Makefile prepare-build-time-extensions

FROM debian:bookworm-slim
ENV DEBIAN_FRONTEND=noninteractive
COPY --from=base-node /usr/bin/make /usr/bin/make
COPY install/docker/Makefile /Makefile
RUN make -f /Makefile install-runtime-deps

COPY ./server/cmd/server-config.env /app/server/cmd/server-config.env
COPY ./server/meshmodel /app/server/meshmodel
COPY ./server/permissions /app/server/permissions/
COPY ./docs/catalog /app/docs/catalog/
COPY --from=meshery-server /meshery /app/server/cmd/
COPY --from=meshery-server /github.com/meshery/meshery/server/helpers/swagger.yaml /app/server/helpers/swagger.yaml
COPY --from=meshery-ui-builder /app/out /app/ui/out
COPY --from=meshery-ui-builder /app/public /app/ui/public
COPY --from=meshery-ui-provider-builder /app/out /app/provider-ui/out
# COPY --from=wrk2 /wrk2 /app/server/cmd/wrk2
# COPY --from=wrk2 /wrk2/wrk /usr/local/bin
COPY --from=seed_content /seed_content /home/appuser/.meshery/seed_content
COPY --from=layer5/getnighthawk:v1.0.1 /usr/local/bin/nighthawk_service /app/server/cmd/
COPY --from=layer5/getnighthawk:v1.0.1 /usr/local/bin/nighthawk_output_transform /app/server/cmd/

# Permission keys are authorization keys which server will
# seed on initial boot for local provider, if there is update on this path
# ensure to update here also https://github.com/meshery/meshery/blob/85a0fd9b9cc03bbc51c4aa1efd1e837db72ef691/install/kubernetes/helm/meshery/values.yaml#L19
ENV KEYS_PATH="/app/server/permissions/keys.csv"

# Build-time Extension Point: Copy .meshery directory from build context if present.
# System integrators can place a .meshery directory at the repository root before running
# docker build to supply build-time configuration, data, packages, etc.
# The contents will be copied into /home/appuser/.meshery/ and merged with existing content.
COPY --from=build_time_extensions /build-ext/.meshery /home/appuser/.meshery

RUN make -f /Makefile create-runtime-user
# Ensure .meshery and subdirectories exist with proper permissions
RUN make -f /Makefile prepare-runtime-dirs
USER appuser
WORKDIR /app/server/cmd
CMD ["./meshery"]
