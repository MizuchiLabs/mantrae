# Builder Stage
FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG TARGETPLATFORM
ENV CGO_ENABLED=1 GOOS=linux GOARCH=${TARGETPLATFORM#linux/}

RUN apk update && apk add --no-cache gcc musl-dev \
  && if [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
  apk add --no-cache gcc-aarch64-linux-gnu; \
  fi

WORKDIR /app
COPY . .

RUN go build -ldflags "-s -w" -o mantrae .

# Final Stage
FROM alpine:latest

RUN apk update && apk add --no-cache bash sqlite

# Copy the binary from the builder stage
COPY --from=builder /app/mantrae /usr/local/bin/mantrae

WORKDIR /app
EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/mantrae"]

