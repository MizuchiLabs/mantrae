#Builder Stage
FROM golang:alpine AS builder

# Install necessary packages for CGO and SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set environment variables for Go
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# Set work directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary with static linking
RUN go build -a -installsuffix cgo -ldflags="-w -s" -o mantrae main.go

# Final Stage
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/mantrae /usr/local/bin/mantrae
COPY entrypoint.sh /

# Ensure the entrypoint script is executable
RUN chmod +x /entrypoint.sh

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 3000

# Define entrypoint and default command
ENTRYPOINT ["/entrypoint.sh"]
CMD ["mantrae"]

