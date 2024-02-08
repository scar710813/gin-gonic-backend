FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o app

EXPOSE 3000

CMD ["./app"]
