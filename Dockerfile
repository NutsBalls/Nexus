# Этап сборки
FROM golang:1.23 AS builder
WORKDIR /app

# Копируем только файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Теперь копируем остальной проект
COPY . .

# Запускаем сборку
RUN go build -o main .

# Этап для запуска
FROM debian:bullseye-slim
WORKDIR /root/

# Копируем скомпилированное приложение из этапа сборки
COPY --from=builder /app/main .

# Копируем конфигурационные файлы, если они есть
COPY --from=builder /app/config /config

# Устанавливаем команду запуска приложения
CMD ["./main"]

