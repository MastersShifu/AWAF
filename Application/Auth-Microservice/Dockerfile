# Первый этап — сборка приложения
FROM golang:1.22.0 AS builder

WORKDIR /app

# Копируем файлы для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и компилируем статический бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Второй этап — минимальный образ для запуска
FROM alpine:latest

WORKDIR /root/

# Добавляем зависимости для корректной работы бинарника и bash
RUN apk add --no-cache ca-certificates bash

# Копируем скомпилированный бинарник из builder-образа
COPY --from=builder /app/main .
COPY ./config.yml /root/config.yml

# Копируем скрипт ожидания
COPY ./wait-for-it.sh /root/wait-for-it.sh

# Делаем скрипт исполнимым
RUN chmod +x /root/wait-for-it.sh

# Открываем порт
EXPOSE 8081

# Запуск приложения с ожиданием запуска базы данных
CMD ["/root/wait-for-it.sh", "auth-postgres:5432", "--", "./main"]
