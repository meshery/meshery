FROM golang:1.16 as meshery-server
ARG TOKEN
ARG GIT_VERSION
ARG GIT_COMMITSHA
ARG RELEASE_CHANNEL

RUN adduser --disabled-login appuser
WORKDIR /github.com/meshery/meshery
ADD . .
RUN go clean -modcache; cd cmd; GOPROXY=https://proxy.golang.org GOSUMDB=off go build -ldflags="-w -s -X main.globalTokenForAnonymousResults=$TOKEN -X main.version=$GIT_VERSION -X main.commitsha=$GIT_COMMITSHA -X main.releasechannel=$RELEASE_CHANNEL" -tags draft -a -o /meshery .

FROM node as ui
ADD ui ui
RUN cd ui; npm install --only=production; npm run build && npm run export; mv out /

FROM node as provider-ui
ADD provider-ui provider-ui
RUN cd provider-ui; npm install --only=production; npm run build && npm run export; mv out /

FROM ubuntu as wrk2
RUN apt-get -y update && apt-get -y install git && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/man/?? /usr/share/man/??_*
RUN apt-get -y update && apt-get -y  install build-essential libssl-dev git zlib1g-dev
RUN git config --global user.email "meshery@layer5.io"
RUN git config --global user.name "meshery"
RUN git clone --depth=1 https://github.com/layer5io/wrk2 && cd wrk2 && make

FROM alpine:3.14 as seed_content
RUN apk add --no-cache curl
WORKDIR /
RUN curl -L -s `curl -s https://api.github.com/repos/layer5io/wasm-filters/releases/latest | grep "browser_download_url" | cut -d : -f 2,3 | tr -d '"'` -o wasm-filters.tar.gz \
    && mkdir -p /seed_content/filters/binaries \
    && tar xzf wasm-filters.tar.gz --directory=/seed_content/filters/binaries

RUN curl -L -s https://github.com/service-mesh-patterns/service-mesh-patterns/tarball/master -o service-mesh-patterns.tgz \
    && mkdir service-mesh-patterns \
    && mkdir -p /seed_content/patterns \
    && tar xzf service-mesh-patterns.tgz --directory=service-mesh-patterns \
    && mv service-mesh-patterns/*/samples/* /seed_content/patterns/

#FROM ubuntu as nighthawk
#RUN apt-get -y update && apt-get -y install git && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/man/?? /usr/share/man/??_*
#RUN apt-get -y update && apt-get -y  install build-essential libssl-dev git zlib1g-dev
#RUN git config --global user.email "meshery@layer5.io"
#RUN git config --global user.name "meshery"
#RUN git clone https://github.com/layer5io/nighthawk-go
#RUN cd nighthawk-go/apinighthawk/bin && chmod +x ./nighthawk_client

FROM alpine:3.14 as jsonschema-util
RUN apk add --no-cache curl
WORKDIR /
RUN UTIL_VERSION=$(curl -L -s https://api.github.com/repos/layer5io/kubeopenapi-jsonschema/releases/latest | \
	grep tag_name | sed "s/ *\"tag_name\": *\"\\(.*\\)\",*/\\1/" | \
	grep -v "rc\.[0-9]$"| head -n 1 ) \
	&& curl -L https://github.com/layer5io/kubeopenapi-jsonschema/releases/download/${UTIL_VERSION}/kubeopenapi-jsonschema-alpine -o kubeopenapi-jsonschema \
	&& chmod +x /kubeopenapi-jsonschema

FROM frolvlad/alpine-glibc:alpine-3.13_glibc-2.32
#RUN apt-get update; apt-get install -y ca-certificates; update-ca-certificates && rm -rf /var/lib/apt/lists/*
RUN apk update && apk add ca-certificates; update-ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates
COPY ./oam /app/oam
COPY --from=meshery-server /meshery /app/cmd/
COPY --from=meshery-server /etc/passwd /etc/passwd
COPY --from=meshery-server /github.com/meshery/meshery/helpers/swagger.yaml /app/helpers/swagger.yaml
COPY --from=ui /out /app/ui/out
COPY --from=provider-ui /out /app/provider-ui/out
COPY --from=wrk2 /wrk2 /app/cmd/wrk2
COPY --from=wrk2 /wrk2/wrk /usr/local/bin
COPY --from=seed_content /seed_content /home/appuser/.meshery/seed_content
COPY --from=layer5/getnighthawk:latest /usr/local/bin/nighthawk_service /app/cmd/
COPY --from=layer5/getnighthawk:latest /usr/local/bin/nighthawk_output_transform /app/cmd/
COPY --from=jsonschema-util /kubeopenapi-jsonschema /home/appuser/.meshery/bin/kubeopenapi-jsonschema

RUN mkdir -p /home/appuser/.meshery/config; chown -R appuser /home/appuser/
USER appuser
WORKDIR /app/cmd
CMD ./meshery
