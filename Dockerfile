# ---------- Build stage ----------
    FROM golang:1.24.2 AS builder

    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
    
    # ---------- Runtime stage ----------
    FROM debian:bookworm-slim
    
    WORKDIR /app
    
    RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
    
    COPY --from=builder /app/server .
    COPY --from=builder /app/docs ./docs
    COPY --from=builder /app/migrations ./migrations
    COPY --from=builder /app/internal/security/keys ./internal/security/keys
    COPY .env .env
    
    EXPOSE 8080
    
    CMD ["./server"]
    