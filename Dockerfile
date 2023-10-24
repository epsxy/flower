# syntax=docker/dockerfile:1

FROM golang:1.21.2

WORKDIR /app

COPY go.mod go.sum ./
COPY Makefile ./Makefile
COPY ./main.go ./main.go
COPY ./VERSION ./VERSION
COPY ./WORKSPACE ./WORKSPACE
COPY ./deps.bzl ./deps.bzl
COPY ./pkg ./pkg
COPY ./bin ./bin

RUN go mod download

#RUN make build-darwin-amd64