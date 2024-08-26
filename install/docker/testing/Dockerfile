FROM golang:1.21.8 as meshery-server
ARG GIT_VERSION

RUN adduser --disabled-login appuser
WORKDIR /github.com/meshery/meshery
COPY go.sum go.mod ./
COPY server/ server/
COPY mesheryctl/ mesheryctl/
RUN go mod download
RUN cd server && cd cmd && GOPROXY=https://proxy.golang.org go build -ldflags="-w -s -X main.version=$GIT_VERSION" -tags draft -a -o /meshery .

# This image will be reuse for both UI & UI Provider
FROM node:18.13.0-alpine3.17 as base-node

FROM base-node as ui-deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

COPY ui/package.json ui/package-lock.json ./
ENV HUSKY=0
RUN npm ci --omit=dev --legacy-peer-deps --ignore-scripts

FROM base-node AS ui-builder
WORKDIR /app
COPY --from=ui-deps /app/node_modules ./node_modules
COPY ui .

# Next.js collects completely anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line in case you want to disable telemetry during the build.
ENV NEXT_TELEMETRY_DISABLED 1

RUN npm run build && npm run export

FROM base-node as ui-provider-deps
RUN apk add --no-cache libc6-compat
WORKDIR /app
# Install dependencies based on the preferred package manager
COPY provider-ui/package.json provider-ui/package-lock.json ./
RUN npm ci --omit=dev --legacy-peer-deps

FROM base-node AS ui-provider-builder
WORKDIR /app
COPY --from=ui-provider-deps /app/node_modules ./node_modules

# Next.js collects completely anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line in case you want to disable telemetry during the build.
ENV NEXT_TELEMETRY_DISABLED 1

COPY provider-ui .
RUN npm run build && npm run export

# This image will be reuse for wrk2, seed content and the last built
FROM alpine:3.15.3 as base-alpine

FROM base-alpine as wrk2
RUN apk update \
    && apk --no-cache --update add alpine-sdk zlib-dev openssl-dev wget git
RUN git config --global user.email "meshery@layer5.io"
RUN git config --global user.name "meshery"
RUN git clone --depth=1 https://github.com/layer5io/wrk2 && cd wrk2 && make

FROM base-alpine as seed_content
RUN apk add --no-cache curl
WORKDIR /
RUN lines=$(curl -s https://api.github.com/repos/layer5io/wasm-filters/releases/latest | grep "browser_download_url.*wasm" | cut -d : -f 2,3 | sed 's/"//g') \
    && mkdir -p seed_content/filters/binaries \
    && cd seed_content/filters/binaries  \
    for line in $lines \
    do \
    curl -LO $line \
    done

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
    && curl -L https://raw.githubusercontent.com/layer5io/image-hub/master/deployment.yaml -o imagehub.yaml \
    && mkdir /emojivoto && (cd /emojivoto && curl --remote-name-all -L https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/emoji.yml \
    https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/vote-bot.yml \
    https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/voting.yml \
    https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/web.yml) \
    && awk 'FNR==1 && NR>1 { printf("\n%s\n\n","---") } 1' /emojivoto/*.yml > /seed_content/applications/emojivoto.yml

#FROM ubuntu as nighthawk
#RUN apt-get -y update && apt-get -y install git && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/man/?? /usr/share/man/??_*
#RUN apt-get -y update && apt-get -y  install build-essential libssl-dev git zlib1g-dev
#RUN git config --global user.email "meshery@layer5.io"
#RUN git config --global user.name "meshery"
#RUN git clone https://github.com/layer5io/nighthawk-go
#RUN cd nighthawk-go/apinighthawk/bin && chmod +x ./nighthawk_client

FROM base-alpine
ENV GLIBC_REPO=https://github.com/sgerrand/alpine-pkg-glibc
ENV GLIBC_VERSION=2.34-r0
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
COPY ./server/meshmodel /app/server/meshmodel
COPY ./server/permissions /app/server/permissions/
COPY --from=meshery-server /meshery /app/server/cmd/
COPY --from=meshery-server /etc/passwd /etc/passwd
COPY --from=meshery-server /github.com/meshery/meshery/server/helpers/swagger.yaml /app/server/helpers/swagger.yaml
COPY --from=ui-builder /app/out /app/ui/out
COPY --from=ui-builder /app/public /app/ui/public
COPY --from=ui-provider-builder /app/out /app/provider-ui/out
COPY --from=wrk2 /wrk2 /app/server/cmd/wrk2
COPY --from=wrk2 /wrk2/wrk /usr/local/bin
COPY --from=seed_content /seed_content /home/appuser/.meshery/seed_content
COPY --from=layer5/getnighthawk:v1.0.1 /usr/local/bin/nighthawk_service /app/server/cmd/
COPY --from=layer5/getnighthawk:v1.0.1 /usr/local/bin/nighthawk_output_transform /app/server/cmd/

RUN addgroup docker
RUN mkdir -p /var/run; chown -R appuser:docker /var/run/
RUN mkdir -p /home/appuser/.meshery/config; chown -R appuser /home/appuser/; chown -R appuser /app/ui/public
USER appuser
WORKDIR /app/server/cmd
CMD ./meshery
