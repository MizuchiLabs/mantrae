FROM gcr.io/distroless/base-debian12:debug

ARG TARGETOS
ARG TARGETARCH

COPY mantrae-${TARGETOS}-${TARGETARCH} /usr/local/bin/mantrae

WORKDIR /app
EXPOSE 3000
EXPOSE 8090

ENTRYPOINT ["/usr/local/bin/mantrae"]

