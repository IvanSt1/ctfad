# backend/Dockerfile

FROM golang:1.23 AS builder
WORKDIR /app

# Копируем go.mod (и, возможно, старый go.sum, если есть)
COPY go.mod ./

# Копируем весь исходный код, чтобы go mod tidy увидел ваши импорты
COPY . .

# Генерируем go.sum и скачиваем зависимости
RUN go mod tidy

# Собираем бинарник с кэшем сборки
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o otkritki

FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/otkritki .
ENTRYPOINT ["./otkritki"]