# Базовый образ с минимальной системой
FROM debian:bullseye-slim AS builder

# Устанавливаем зависимости
RUN apt-get update && apt-get install -y \
    curl wget unzip git libglu1-mesa libxi6 libxrender1 xz-utils apt-transport-https ca-certificates && \
    apt-get clean

# Создаем пользователя для Flutter
RUN useradd -ms /bin/bash flutteruser

# Устанавливаем Flutter SDK
WORKDIR /usr/local
RUN git clone https://github.com/flutter/flutter.git -b stable --depth 1

# Добавляем Flutter в PATH
ENV PATH="/usr/local/flutter/bin:$PATH"

# Меняем владельца Flutter SDK и временной директории
RUN chown -R flutteruser:flutteruser /usr/local/flutter && \
    mkdir -p /tmp/flutter_temp && \
    chown -R flutteruser:flutteruser /tmp/flutter_temp && \
    chmod -R 777 /tmp/flutter_temp

# Устанавливаем переменные окружения для Flutter
ENV PUB_CACHE="/tmp/flutter_temp/.pub-cache"
ENV PUB_HOSTED_URL="https://pub.dev"
ENV FLUTTER_STORAGE_BASE_URL="https://storage.googleapis.com"
ENV TMPDIR="/tmp/flutter_temp"

# Переход на пользователя flutteruser
USER flutteruser
WORKDIR /usr/local/flutter

# Проверяем установку Flutter
RUN flutter doctor

# Копируем проект с назначением владельца
WORKDIR /app
COPY --chown=flutteruser:flutteruser mobile_app/ . 

# Устанавливаем зависимости
RUN flutter pub get

# Собираем проект
RUN flutter build web --release

# Минимальный веб-сервер для запуска Flutter Web
FROM nginx:alpine
WORKDIR /usr/share/nginx/html
COPY --from=builder /app/build/web . 
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
