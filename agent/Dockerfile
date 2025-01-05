FROM alpine

ARG VERSION
ARG COMMIT
ARG DATE

# Metadata labels
LABEL org.opencontainers.image.vendor="Mizuchi Labs"
LABEL org.opencontainers.image.source="https://github.com/MizuchiLabs/mantrae"
LABEL org.opencontainers.image.title="Mantrae Agent"
LABEL org.opencontainers.image.version=$VERSION
LABEL org.opencontainers.image.revision=$COMMIT
LABEL org.opencontainers.image.created=$DATE
LABEL org.opencontainers.image.licenses="MIT"

WORKDIR /app

COPY mantrae_agent /usr/local/bin/mantrae_agent

EXPOSE 3000
EXPOSE 8090

ENTRYPOINT ["/usr/local/bin/mantrae_agent"]

