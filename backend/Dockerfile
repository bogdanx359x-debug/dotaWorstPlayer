FROM golang:1.24-alpine AS build

WORKDIR /app

# Копіюємо файли з залежностями
COPY go.mod go.sum ./
RUN go mod download

# Копіюємо решту коду
COPY . .

# Будуємо бінарник (головний файл лежить у internal/)
RUN go build -o server ./internal

# Мінімальний образ для запуску
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]
