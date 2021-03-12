FROM golang:1.13.6 as meshery-server
ARG TOKEN
ARG GIT_VERSION
ARG GIT_COMMITSHA
ARG RELEASE_CHANNEL

RUN adduser --disabled-login appuser
WORKDIR /github.com/layer5io/meshery
ADD . .
RUN rm go.sum; go clean -modcache; cd cmd; GOPROXY=https://proxy.golang.org GOSUMDB=off go build -ldflags="-w -s -X main.globalTokenForAnonymousResults=$TOKEN -X main.version=$GIT_VERSION -X main.commitsha=$GIT_COMMITSHA -X main.releasechannel=$RELEASE_CHANNEL" -tags draft -a -o /meshery .

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

FROM ubuntu as nighthawk
RUN apt-get -y update && apt-get -y install git && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/man/?? /usr/share/man/??_*
RUN apt-get -y update && apt-get -y  install build-essential libssl-dev git zlib1g-dev
RUN git config --global user.email "meshery@layer5.io"
RUN git config --global user.name "meshery"
RUN git clone https://github.com/layer5io/nighthawk-go
RUN cd nighthawk-go/apinighthawk/bin && chmod +x ./nighthawk_client

FROM ubuntu
RUN apt-get update; apt-get install -y ca-certificates; update-ca-certificates
COPY ./oam /app/oam
COPY --from=meshery-server /meshery /app/cmd/
COPY --from=meshery-server /etc/passwd /etc/passwd
COPY --from=ui /out /app/ui/out
COPY --from=provider-ui /out /app/provider-ui/out
COPY --from=wrk2 /wrk2 /app/cmd/wrk2
COPY --from=wrk2 /wrk2/wrk /usr/local/bin
COPY --from=nighthawk /nighthawk-go/apinighthawk/bin /usr/local/bin
RUN mkdir -p /home/appuser/.meshery/config; chown -R appuser /home/appuser/
USER appuser
WORKDIR /app/cmd
CMD ./meshery
