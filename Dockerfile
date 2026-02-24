# Этап 1: Сборка
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum (если есть)
COPY go.mod go.sum* ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hot-coffee .

# Этап 2: Финальный образ
FROM alpine:latest

# Устанавливаем CA сертификаты (для HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/hot-coffee .

# Копируем папку data со стартовыми данными
COPY --from=builder /app/data ./data
COPY --from=builder /app/frontend ./frontend
# Создаём директорию data, если её нет
RUN mkdir -p /root/data

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./hot-coffee"]
