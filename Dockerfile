FROM golang:1.12.6 as bd0
RUN adduser --disabled-login appuser
WORKDIR /github.com/layer5io/meshery
ADD . .
RUN cd cmd; go build -ldflags="-w -s" -tags draft -a -o /meshery .

FROM node as bd1
ADD ui ui
RUN cd ui; npm i; npm run build && npm run export; mv out /

# FROM alpine
# RUN apk --update add ca-certificates
# RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
FROM ubuntu
RUN apt-get update; apt-get install -y ca-certificates; update-ca-certificates
COPY --from=bd0 /meshery /app/cmd/
COPY --from=bd0 /etc/passwd /etc/passwd
COPY --from=bd1 /out /app/ui/out
USER appuser
WORKDIR /app/cmd
CMD ./meshery
