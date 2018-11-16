FROM golang:1.11.2 as bd
WORKDIR /app
ADD . .
RUN go build -a -o /istio-playground .

FROM alpine
RUN apk --update add ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=bd /istio-playground /app/
ADD dashboard.html /app/
ADD get-ao-token.html /app/
ADD static /app/static
ADD templates /app/templates
WORKDIR /app/
CMD ./istio-playground
