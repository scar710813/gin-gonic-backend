FROM golang:1.21.6-alpine as builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go build -o main

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 3001
CMD ["./main"]
