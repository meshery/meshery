FROM golang:1.19 as meshery-server
ARG TOKEN
ARG GIT_VERSION
ARG GIT_COMMITSHA
ARG RELEASE_CHANNEL

RUN adduser --disabled-login appuser
WORKDIR /github.com/meshery/meshery
ADD . .
RUN go clean -modcache; cd server; cd cmd; GOPROXY=https://proxy.golang.org GOSUMDB=off go build -ldflags="-w -s -X main.globalTokenForAnonymousResults=$TOKEN -X main.version=$GIT_VERSION -X main.commitsha=$GIT_COMMITSHA -X main.releasechannel=$RELEASE_CHANNEL" -tags draft -a -o /meshery .

FROM node:lts-slim as ui
ADD ui ui
RUN cd ui; npm install --only=production; npm run build && npm run export; mv out /

FROM node:16 as provider-ui
ADD provider-ui provider-ui
RUN cd provider-ui; npm install --only=production; npm run build && npm run export; mv out /

FROM ubuntu as wrk2
RUN apt-get -y update && apt-get -y install git && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/man/?? /usr/share/man/??_*
RUN apt-get -y update && apt-get -y  install build-essential libssl-dev git zlib1g-dev
RUN git config --global user.email "meshery@layer5.io"
RUN git config --global user.name "meshery"
RUN git clone --depth=1 https://github.com/layer5io/wrk2 && cd wrk2 && make

FROM alpine:3.15.3 as seed_content
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


FROM frolvlad/alpine-glibc:alpine-3.13_glibc-2.32
#RUN apt-get update; apt-get install -y ca-certificates; update-ca-certificates && rm -rf /var/lib/apt/lists/*
RUN apk update && apk add ca-certificates; update-ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates
RUN apk upgrade --no-cache && \
    apk add --no-cache libstdc++
COPY ./server/oam /app/server/oam
COPY --from=meshery-server /meshery /app/server/cmd/
COPY --from=meshery-server /etc/passwd /etc/passwd
COPY --from=meshery-server /github.com/meshery/meshery/server/helpers/swagger.yaml /app/server/helpers/swagger.yaml
COPY --from=ui /out /app/ui/out
COPY --from=provider-ui /out /app/provider-ui/out
COPY --from=wrk2 /wrk2 /app/server/cmd/wrk2
COPY --from=wrk2 /wrk2/wrk /usr/local/bin
COPY --from=seed_content /seed_content /home/appuser/.meshery/seed_content
COPY --from=layer5/getnighthawk:latest /usr/local/bin/nighthawk_service /app/server/cmd/
COPY --from=layer5/getnighthawk:latest /usr/local/bin/nighthawk_output_transform /app/server/cmd/

RUN mkdir -p /home/appuser/.meshery/config; chown -R appuser /home/appuser/
USER appuser
WORKDIR /app/server/cmd
CMD ./meshery
