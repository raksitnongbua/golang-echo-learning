FROM golang:1.19alpine:3.16
WORKDIR /app
COPY . /app

RUN go build -o main main.go
