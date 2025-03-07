FROM golang:1.21-alpine AS builder

# Установка зависимостей системы
RUN apk add --no-cache ca-certificates gcc libc-dev

# Установка рабочей директории
WORKDIR /app

# Копирование исходного кода
COPY . .

# Инициализация модуля и сборка приложения 
RUN go mod tidy && CGO_ENABLED=0 go build -ldflags="-s -w" -o summvideo-app .

# Второй этап сборки
FROM alpine:3.19

# Установка ffmpeg
RUN apk add --no-cache ffmpeg

# Копирование бинарного файла из первого этапа
COPY --from=builder /app/summvideo-app /usr/local/bin/summvideo-app

# Создание директорий для монтирования
RUN mkdir -p /app/static /app/templates

# Установка рабочей директории
WORKDIR /app

# Настройка переменных окружения
ENV GIN_MODE=release
ENV PORT=8000

# Создание директории для логов
RUN mkdir -p /var/log/summvideo

# Открытие порта
EXPOSE 8000

# Запуск приложения
CMD ["/usr/local/bin/summvideo-app"]