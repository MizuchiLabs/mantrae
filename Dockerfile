FROM cgr.dev/chainguard/static:latest
WORKDIR /app
COPY bin/mantrae web/build ./
CMD ["./mantrae"]
