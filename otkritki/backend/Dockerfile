# backend/Dockerfile

# 1) Use Go 1.23 for the build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем go.mod и go.sum и скачиваем модули
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код и собираем бинарь
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o otkritki

# 2) Минимальный рантайм-образ
FROM gcr.io/distroless/base-debian11

WORKDIR /app
COPY --from=builder /app/otkritki .

# Точка входа
ENTRYPOINT ["./otkritki"]

