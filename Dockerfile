# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

ENV GO_ENV=product

RUN apk update && apk add git
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build

EXPOSE 80

CMD ["go", "run", "main.go"]