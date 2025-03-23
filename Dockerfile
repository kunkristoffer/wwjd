ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the Go app
RUN go build -v -o /run-app .

# Final stage
FROM debian:bookworm

# Install build dependencies (optional)
RUN apt-get update && apt-get install -y ca-certificates fuse3 sqlite3 && update-ca-certificates

COPY --from=builder /run-app /usr/local/bin/

COPY --from=builder /usr/src/app/assets ./assets

CMD ["run-app"]
