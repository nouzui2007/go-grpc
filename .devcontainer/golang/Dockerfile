FROM golang:1.14.4-alpine3.12

ENV GO111MODULE=on
ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin
ENV CGO_ENABLED=0

WORKDIR /var/dsaf-api

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64

RUN apk update && \
    apk add git && \
    apk add mysql-client && \
    apk add mysql-dev && \
    go get golang.org/x/tools/cmd/godoc &&\
    rm -rf /var/cache/apk/*
