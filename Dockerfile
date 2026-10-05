# ===== Build stage =====
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Dependency cache: copy module files first.
# If go.mod/go.sum don't change, Docker reuses this layer.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code and build.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# ===== Runtime stage =====
FROM alpine:3.20

# ca-certificates is required for outbound HTTPS requests.
# tzdata ensures time.Now() respects timezones (we use UTC, but it's good practice).
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -u 1000 appuser

WORKDIR /app

# Create /data with appuser ownership BEFORE declaring VOLUME.
# This ensures the initialized volume inherits correct permissions.
RUN mkdir -p /data && chown appuser:appuser /data

USER appuser

# Copy only the binary from the previous stage — no source code, no toolchain.
COPY --from=builder /out/api .

EXPOSE 8080
ENV PORT=8080
ENV DB_PATH=/data/pismo.db

VOLUME ["/data"]

ENTRYPOINT ["./api"]