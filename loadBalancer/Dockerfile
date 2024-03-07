FROM golang:1.21.6-alpine as builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

RUN go build -o shinsei main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/shinsei .

CMD ["./shinsei"]
