FROM golang:1.13.4 as bd0
RUN adduser --disabled-login appuser
WORKDIR /github.com/layer5io/meshery
ADD . .
RUN cd cmd; GOPROXY=direct GOSUMDB=off go build -ldflags="-w -s" -tags draft -a -o /meshery .

FROM node as bd1
ADD ui ui
RUN cd ui; npm i; npm run build && npm run export; mv out /

FROM ubuntu as bd2
RUN apt-get -y update && apt-get -y install git && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/man/?? /usr/share/man/??_*
RUN apt-get -y update && apt-get -y  install build-essential libssl-dev git zlib1g-dev
RUN git config --global user.email "meshery@layer5.io"
RUN git config --global user.name "meshery"
RUN git clone --depth=1 https://github.com/layer5io/wrk2 && cd wrk2 && make 

# FROM alpine
# RUN apk --update add ca-certificates
# RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
FROM ubuntu
RUN apt-get update; apt-get install -y ca-certificates; update-ca-certificates
COPY --from=bd0 /meshery /app/cmd/
COPY --from=bd0 /etc/passwd /etc/passwd
COPY --from=bd1 /out /app/ui/out
COPY --from=bd2 /wrk2 /app/cmd/wrk2
COPY --from=bd2 /wrk2/wrk /usr/local/bin
RUN mkdir -p /home/appuser/.meshery/config; chown -R appuser /home/appuser/
USER appuser
WORKDIR /app/cmd
CMD ./meshery
