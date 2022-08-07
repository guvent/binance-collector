FROM golang:1.19.0-alpine3.15

ENV GO111MODULE=on
ENV GOPATH=""

RUN apk add --no-cache git g++ make libffi-dev librdkafka-dev pkgconf


