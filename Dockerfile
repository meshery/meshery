FROM frolvlad/alpine-glibc:alpine-3.13_glibc-2.32

COPY --from=layer5/getnighthawk:latest /usr/local/bin/nighthawk_service /app/server/cmd/
COPY --from=layer5/getnighthawk:latest /usr/local/bin/nighthawk_output_transform /app/server/cmd/
ENTRYPOINT ["/bin/sh","-c","sleep infinity"]
