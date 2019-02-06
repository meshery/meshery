FROM golang:1.11.5 as bd
RUN adduser --disabled-login appuser
WORKDIR /github.com/layer5io/meshery
ADD . .
RUN cd cmd; go build -ldflags="-w -s" -a -o /meshery .
RUN find . -name "*.go" -type f -delete; mv public /; mv meshes /
RUN cd /public/static/js; wget https://raw.githubusercontent.com/fortio/fortio/master/ui/static/js/Chart.min.js
RUN cd /public/static/js; wget https://raw.githubusercontent.com/fortio/fortio/master/ui/static/js/fortio_chart.js

FROM alpine
RUN apk --update add ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=bd /meshery /app/cmd/
COPY --from=bd /public /app/public
COPY --from=bd /meshes /app/meshes
COPY --from=bd /etc/passwd /etc/passwd
USER appuser
WORKDIR /app/cmd
CMD ./meshery
