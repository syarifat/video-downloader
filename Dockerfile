# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for go modules)
RUN apk add --no-cache git

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o video-bot .

# Stage 2: Runtime image
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \
    python3 \
    py3-pip \
    && pip3 install --break-system-packages yt-dlp \
    && mkdir -p /app/downloads

# Copy the binary from builder
COPY --from=builder /app/video-bot .

# Set default environment variables
ENV DOWNLOAD_DIR=/app/downloads
ENV MAX_FILE_SIZE_MB=50

# Run the bot
CMD ["./video-bot"]
