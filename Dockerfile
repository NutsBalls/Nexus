# Этап сборки
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

RUN mkdir -p /app/uploads && chmod 777 /app/uploads
# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Собираем приложение
RUN go build -o main .

# Этап создания минимального контейнера
FROM debian:bullseye-slim

FROM frolvlad/alpine-glibc:alpine-3.18

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем приложение
COPY --from=builder /app/main .

# Устанавливаем порт
EXPOSE 8080

# Устанавливаем команду по умолчанию
CMD ["./main"]