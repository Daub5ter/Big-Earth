FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY .. /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o parsing-service ./cmd

RUN chmod +x /app/parsing-service

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/parsing-service /app

CMD [ "app/parsing-service" ]



