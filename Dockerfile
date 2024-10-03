# Builder Stage
FROM golang:latest AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=1 go build -ldflags "-s -w" -o mantrae .

# Final Stage
FROM gcr.io/distroless/base-debian12:debug

# Copy the binary from the builder stage
COPY --from=builder /app/mantrae /usr/local/bin/mantrae

WORKDIR /app
EXPOSE 3000
EXPOSE 8090

ENTRYPOINT ["/usr/local/bin/mantrae"]

