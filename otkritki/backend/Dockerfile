# backend/Dockerfile

# 1) Builder stage: compile static Go binary
FROM golang:1.23 AS builder

# Disable CGO for static binary
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy go.mod first
COPY go.mod ./

# Copy all source files
COPY . .

# Generate go.sum and download all dependencies
RUN go mod tidy \
    && go mod download

# Build static binary with optimizations
RUN go build -o otkritki -ldflags="-s -w"

# 2) Final stage: minimal image with statically linked binary
FROM gcr.io/distroless/static:nonroot AS runner

WORKDIR /app
COPY --from=builder /app/otkritki ./otkritki
USER nonroot:nonroot
ENTRYPOINT ["./otkritki"]: minimal image with statically linked binary
FROM gcr.io/distroless/static:nonroot AS runner

WORKDIR /app
COPY --from=builder /app/otkritki ./otkritki
USER nonroot:nonroot
ENTRYPOINT ["./otkritki"]
