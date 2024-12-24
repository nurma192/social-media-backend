# Обновляем базовый образ до последней версии Go
FROM golang:alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код проекта в контейнер
COPY . .

# Устанавливаем переменные окружения для PostgreSQL
ENV DB_HOST=db \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=uk888888 \
    DB_NAME=sm


# Компилируем приложение
RUN go build -o main ./cmd/main.go

# Указываем порт, который будет использован приложением
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]

# Компилируем приложение
RUN go build -o main ./cmd/main.go

# Указываем порт, который будет использован приложением
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
