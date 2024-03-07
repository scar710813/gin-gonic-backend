FROM golang:1.21.6-alpine as builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

RUN go build -o main

FROM alpine:latest

COPY --from=builder /app/main /app/main

EXPOSE 3000

CMD ["/app/main"]
