# Dockerfile
# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем зависимости первыми для лучшего кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем ВСЕ исходные файлы (включая internal и cmd)
COPY . .

# Собираем бинарник с указанием правильного пути
RUN CGO_ENABLED=0 GOOS=linux go build -o /bot ./cmd/bot/main.go

# Final stage
FROM alpine:3.18

WORKDIR /app

# Копируем бинарник и необходимые файлы
COPY --from=builder /bot /app/bot
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/tarantool_init.lua /app/tarantool_init.lua

# Устанавливаем зависимости для TLS
RUN apk --no-cache add ca-certificates

CMD ["/app/bot"]