FROM golang:1.22 as builder

WORKDIR /app
COPY . .

RUN go mod download
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o otkritki

FROM ubuntu:22.04
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/otkritki .
ENTRYPOINT ["./otkritki"]
