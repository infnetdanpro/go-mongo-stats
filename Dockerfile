FROM golang:1.18

WORKDIR /www/

COPY go.mod .
COPY go.sum .

RUN apt update && apt install gcc
ENV CGO_ENABLED 1
RUN go mod download

EXPOSE 8088