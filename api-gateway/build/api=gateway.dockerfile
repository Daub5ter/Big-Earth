FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY .. /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o api-gateway ./cmd

RUN chmod +x /app/api-gateway

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/api-gateway /app

CMD [ "app/api-gateway" ]