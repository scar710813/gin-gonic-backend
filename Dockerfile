FROM golang:1.21.6-alpine as builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

RUN go build -o main main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/healthcheck.sh .
COPY --from=builder /app/schema.sql .

EXPOSE 3000

CMD ["./main"]
