FROM cgr.dev/chainguard/static:latest
WORKDIR /app
COPY mantrae .
EXPOSE 3000
CMD ["./mantrae"]
