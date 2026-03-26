# Используем официальный образ Go
FROM golang:1.23-alpine

# Рабочая директория внутри контейнера
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/app/main.go

# Запускаем
CMD ["./main"]