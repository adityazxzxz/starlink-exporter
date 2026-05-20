# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build arguments
ARG VERSION=dev
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown
ARG GO_VERSION=unknown

# Build
RUN CGO_ENABLED=0 go build \
  -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -X main.GoVersion=${GO_VERSION}" \
  -o starlink_exporter ./cmd/starlink_exporter

# Runtime stage
FROM gcr.io/distroless/static:nonroot

WORKDIR /app
USER nonroot:nonroot

# Copy binary
COPY --from=builder /app/starlink_exporter /app/starlink_exporter

# Default environment variables
ENV MODE=web \
    SOURCE=dummy \
    LISTEN=:9817 \
    ADDRESS=192.168.100.1:9200 \
    PUSHGATEWAY="" \
    JOB=starlink_exporter \
    INSTANCE=starlink \
    INTERVAL=15s \
    LOG_LEVEL=info

# Expose port
EXPOSE 9817

# Entry point that constructs command from environment variables
ENTRYPOINT ["/app/starlink_exporter", \
            "-mode=${MODE}", \
            "-source=${SOURCE}", \
            "-listen=${LISTEN}", \
            "-address=${ADDRESS}", \
            "-pushgateway=${PUSHGATEWAY}", \
            "-job=${JOB}", \
            "-instance=${INSTANCE}", \
            "-interval=${INTERVAL}", \
            "-log-level=${LOG_LEVEL}"]
