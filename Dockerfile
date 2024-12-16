
FROM golang:1.23 AS builder
WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go run main.go


FROM debian:bullseye-slim
WORKDIR /root/


COPY --from=builder /app/main .


COPY --from=builder /app/config /config


CMD ["./main"]

