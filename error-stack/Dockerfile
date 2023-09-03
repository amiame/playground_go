FROM golang:1.20.7-alpine AS base

WORKDIR /go/src/proj
COPY ./ .

RUN ["go", "build", "-o", "app", "main/main.go"]

ENTRYPOINT ["./app"]
