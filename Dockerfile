FROM golang:1.23.2-bookworm

WORKDIR /app

COPY . /app

EXPOSE 8080