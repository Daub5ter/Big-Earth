FROM golang:1.21.1-alpine as builder

RUN mkdir /app

COPY .. /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o parsing-service ./cmd

RUN chmod +x /app/parsing-service

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/parsing-service /app

COPY /build/cert.pem /app/cert.pem

COPY /build/key.pem /app/key.pem

COPY /configs/parsing-config.yaml /app/parsing-config.yaml

CMD [ "app/parsing-service" ]