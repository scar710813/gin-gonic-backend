FROM golang:alpine

WORKDIR /app

COPY . .

COPY database/nginx.conf /etc/nginx/nginx.conf
COPY database/schema.sql docker-entrypoint-initdb.d/schema.sql

RUN go build -o rinha_kk

EXPOSE 3000

CMD ["./rinha_kk"]
