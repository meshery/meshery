FROM golang:1.25.5 AS meshery-server
ARG GIT_VERSION

WORKDIR /github.com/meshery/meshery
COPY go.sum go.mod ./
COPY server/ server/
COPY mesheryctl/ mesheryctl/
COPY install/docker/testing/Makefile /github.com/meshery/meshery/install/docker/testing/Makefile
RUN make -f install/docker/testing/Makefile download-go-modules
RUN make -f install/docker/testing/Makefile build-meshery-server GIT_VERSION="$GIT_VERSION"

# This image will be reuse for both UI & UI Provider
FROM node:20 AS base-node

FROM base-node AS ui-deps
WORKDIR /app
COPY ui/package.json ui/package-lock.json ./
COPY install/docker/testing/Makefile /app/Makefile
ENV HUSKY=0
RUN make -f /app/Makefile ui-deps

FROM base-node AS ui-builder
WORKDIR /app
COPY --from=ui-deps /app/node_modules ./node_modules
COPY install/docker/testing/Makefile /app/Makefile
COPY ui .

# Next.js collects completely anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line in case you want to disable telemetry during the build.
ENV NEXT_TELEMETRY_DISABLED=1

RUN make -f /app/Makefile ui-build

FROM base-node AS ui-provider-deps
WORKDIR /app
# Install dependencies based on the preferred package manager
COPY provider-ui/package.json provider-ui/package-lock.json ./
COPY install/docker/testing/Makefile /app/Makefile
RUN make -f /app/Makefile provider-ui-deps

FROM base-node AS ui-provider-builder
WORKDIR /app
COPY --from=ui-provider-deps /app/node_modules ./node_modules
COPY install/docker/testing/Makefile /app/Makefile

# Next.js collects completely anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line in case you want to disable telemetry during the build.
ENV NEXT_TELEMETRY_DISABLED=1

COPY provider-ui .
RUN make -f /app/Makefile provider-ui-build

FROM golang:1.25.5 AS wrk2
COPY install/docker/testing/Makefile /Makefile
RUN make -f /Makefile install-wrk2-build-deps
RUN make -f /Makefile configure-wrk2-git
RUN make -f /Makefile clone-and-build-wrk2

FROM debian:bookworm-slim AS seed_content
COPY --from=base-node /usr/bin/make /usr/bin/make
COPY install/docker/testing/Makefile /Makefile
WORKDIR /
RUN make -f /Makefile install-seed-content-deps
RUN make -f /Makefile seed-filters

# bundling patterns
RUN make -f /Makefile seed-patterns

# bundling applications
RUN make -f /Makefile seed-applications

#FROM ubuntu as nighthawk
#RUN apt-get -y update && apt-get -y install git && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/man/?? /usr/share/man/??_*
#RUN apt-get -y update && apt-get -y  install build-essential libssl-dev git zlib1g-dev
#RUN git config --global user.email "meshery@meshery.io"
#RUN git config --global user.name "meshery"
#RUN git clone https://github.com/layer5io/nighthawk-go
#RUN cd nighthawk-go/apinighthawk/bin && chmod +x ./nighthawk_client

FROM debian:bookworm-slim
ENV DEBIAN_FRONTEND=noninteractive
COPY --from=base-node /usr/bin/make /usr/bin/make
COPY install/docker/testing/Makefile /Makefile
RUN make -f /Makefile install-runtime-deps

# FROM frolvlad/alpine-glibc:alpine-3.13_glibc-2.32
#RUN apt-get update; apt-get install -y ca-certificates; update-ca-certificates && rm -rf /var/lib/apt/lists/*
# RUN apk update && apk add ca-certificates; update-ca-certificates && rm -rf /var/cache/apk/*
# RUN apk upgrade --no-cache && \
#     apk add --no-cache libstdc++

COPY ./server/cmd/server-config.env /app/server/cmd/server-config.env
COPY ./server/meshmodel /app/server/meshmodel
COPY ./server/permissions /app/server/permissions/
COPY --from=meshery-server /meshery /app/server/cmd/
COPY --from=meshery-server /github.com/meshery/meshery/server/helpers/swagger.yaml /app/server/helpers/swagger.yaml
COPY --from=ui-builder /app/out /app/ui/out
COPY --from=ui-builder /app/public /app/ui/public
COPY --from=ui-provider-builder /app/out /app/provider-ui/out
COPY --from=wrk2 /wrk2 /app/server/cmd/wrk2
COPY --from=wrk2 /wrk2/wrk /usr/local/bin
COPY --from=seed_content /seed_content /home/appuser/.meshery/seed_content
COPY --from=layer5/getnighthawk:v1.0.1 /usr/local/bin/nighthawk_service /app/server/cmd/
COPY --from=layer5/getnighthawk:v1.0.1 /usr/local/bin/nighthawk_output_transform /app/server/cmd/

RUN make -f /Makefile create-runtime-user
RUN make -f /Makefile prepare-runtime-dirs
USER appuser
WORKDIR /app/server/cmd
CMD ["./meshery"]
