FROM golang:1.9
# Base image that already has the pre-requisites downloaded.

WORKDIR /go/src/github.com/zmap

RUN go-wrapper download github.com/WangYihang/zgrab2

WORKDIR /go/src/github.com/WangYihang/zgrab2

RUN go get -v ./...
RUN go get -v -t ./...
