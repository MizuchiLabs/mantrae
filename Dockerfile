# Builder Stage
FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG TARGETPLATFORM
ENV CGO_ENABLED=1 

RUN xx-apk add musl-dev gcc

WORKDIR /app
COPY . .

RUN xx-go build -ldflags "-s -w" -o mantrae .

# Final Stage
FROM alpine:latest

RUN apk update && apk add --no-cache bash sqlite

# Copy the binary from the builder stage
COPY --from=builder /app/mantrae /usr/local/bin/mantrae

WORKDIR /app
EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/mantrae"]

