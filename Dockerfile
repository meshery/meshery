FROM frolvlad/alpine-glibc:alpine-3.13_glibc-2.32@sha256:cc9a97ed4e27fe0129056fa422d1c54f5a2f7ebe9a2856d7d5a15f17a17614fc

COPY --from=layer5/getnighthawk:latest /usr/local/bin/nighthawk_service /app/server/cmd/
COPY --from=layer5/getnighthawk:latest /usr/local/bin/nighthawk_output_transform /app/server/cmd/
ENTRYPOINT ["/bin/sh","-c","sleep infinity"]
