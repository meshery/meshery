FROM golang:1.26.4 AS meshery-server
ARG GIT_VERSION

RUN adduser --disabled-login appuser
WORKDIR /github.com/meshery/meshery
COPY go.sum go.mod ./
COPY server/ server/
COPY mesheryctl/ mesheryctl/
RUN go mod download
RUN cd server && cd cmd && GOPROXY=https://proxy.golang.org go build -ldflags="-w -s -X main.version=$GIT_VERSION" -tags draft -a -o /meshery .

# This image will be reuse for both UI & UI Provider
FROM node:22-alpine AS base-node

FROM base-node AS ui-deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

COPY ui/package.json ui/package-lock.json ./
ENV HUSKY=0
RUN npm ci --legacy-peer-deps --ignore-scripts

FROM base-node AS ui-builder
WORKDIR /app
COPY --from=ui-deps /app/node_modules ./node_modules
COPY ui .

# Next.js collects completely anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line in case you want to disable telemetry during the build.
ENV NEXT_TELEMETRY_DISABLED=1

RUN npm run build

FROM base-node AS ui-provider-deps
RUN apk add --no-cache libc6-compat
WORKDIR /app
# Install dependencies based on the preferred package manager
COPY provider-ui/package.json provider-ui/package-lock.json ./
RUN npm ci --legacy-peer-deps

FROM base-node AS ui-provider-builder
WORKDIR /app
COPY --from=ui-provider-deps /app/node_modules ./node_modules

# Next.js collects completely anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line in case you want to disable telemetry during the build.
ENV NEXT_TELEMETRY_DISABLED=1

COPY provider-ui .
RUN npm run build

# This image is reused for seed content and the last built
FROM alpine:3.23.2 AS base-alpine

FROM base-alpine AS seed_content
RUN apk add --no-cache curl
WORKDIR /

# bundling patterns
RUN curl -L -s https://github.com/service-mesh-patterns/service-mesh-patterns/tarball/master -o service-mesh-patterns.tgz \
    && mkdir service-mesh-patterns \
    && mkdir -p /seed_content/patterns \
    && tar xzf service-mesh-patterns.tgz --directory=service-mesh-patterns \
    && mv service-mesh-patterns/*/samples/* /seed_content/patterns/

# bundling applications
RUN mkdir -p /seed_content/applications && cd /seed_content/applications \
    && curl -LO https://raw.githubusercontent.com/istio/istio/master/samples/bookinfo/platform/kube/bookinfo.yaml \
    && curl -LO https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml \
    && mkdir /emojivoto && (cd /emojivoto && curl --remote-name-all -L https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/emoji.yml \
    https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/vote-bot.yml \
    https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/voting.yml \
    https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/web.yml) \
    && awk 'FNR==1 && NR>1 { printf("\n%s\n\n","---") } 1' /emojivoto/*.yml > /seed_content/applications/emojivoto.yml

FROM base-alpine
ENV GLIBC_REPO=https://github.com/sgerrand/alpine-pkg-glibc
ENV GLIBC_VERSION=2.35-r1
# To solve: ./meshery: /usr/glibc-compat/lib/libc.so.6: version `GLIBC_2.34' not found (required by ./meshery)
RUN set -ex && \
    apk --update add libstdc++ curl ca-certificates && \
    for pkg in glibc-${GLIBC_VERSION} glibc-bin-${GLIBC_VERSION}; \
    do curl -sSL ${GLIBC_REPO}/releases/download/${GLIBC_VERSION}/${pkg}.apk -o /tmp/${pkg}.apk; done && \
    apk add --allow-untrusted /tmp/*.apk && \
    rm -v /tmp/*.apk && \
    /usr/glibc-compat/sbin/ldconfig /lib /usr/glibc-compat/lib
RUN update-ca-certificates

# FROM frolvlad/alpine-glibc:alpine-3.13_glibc-2.32
#RUN apt-get update; apt-get install -y ca-certificates; update-ca-certificates && rm -rf /var/lib/apt/lists/*
# RUN apk update && apk add ca-certificates; update-ca-certificates && rm -rf /var/cache/apk/*
# RUN apk upgrade --no-cache && \
#     apk add --no-cache libstdc++

COPY ./server/cmd/server-config.env /app/server/cmd/server-config.env
COPY ./models /app/models
COPY ./server/permissions /app/server/permissions/
COPY --from=meshery-server /meshery /app/server/cmd/
COPY --from=meshery-server /etc/passwd /etc/passwd
COPY --from=ui-builder /app/out /app/ui/out
COPY --from=ui-builder /app/public /app/ui/public
COPY --from=ui-provider-builder /app/out /app/provider-ui/out
COPY --from=seed_content /seed_content /home/appuser/.meshery/seed_content

RUN addgroup docker
RUN mkdir -p /var/run; chown -R appuser:docker /var/run/
RUN mkdir -p /home/appuser/.meshery/config; chown -R appuser /home/appuser/; chown -R appuser /app/ui/public
USER appuser
WORKDIR /app/server/cmd
CMD ["./meshery"]
