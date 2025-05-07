FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Статическая сборка
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM scratch

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/docs ./docs
COPY .env .env
COPY internal/security/keys ./internal/security/keys

EXPOSE 8080
CMD ["./server"]
