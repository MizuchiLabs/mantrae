# Builder Stage
FROM golang:alpine AS builder
ENV CGO_ENABLED=1

RUN apk update && apk add --no-cache gcc musl-dev

WORKDIR /app
COPY . .

RUN go build -o mantrae .

# Final Stage
FROM alpine:latest

RUN apk update && apk add --no-cache bash sqlite

# Copy the binary from the builder stage
COPY --from=builder /app/mantrae /usr/local/bin/mantrae
COPY entrypoint.sh /

RUN chmod +x /entrypoint.sh

WORKDIR /app
EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/mantrae"]

