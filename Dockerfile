FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o rinha_kk

EXPOSE 3000

CMD ["./rinha_kk"]
