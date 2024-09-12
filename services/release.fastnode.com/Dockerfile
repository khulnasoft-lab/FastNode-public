FROM golang:1.15.3-buster as builder

ENV GOPRIVATE=github.com/khulnasoft-lab/*

WORKDIR /go/src/github.com/khulnasoft-lab/fastnode

COPY go.mod .
COPY go.sum .
COPY fastnode-go ./fastnode-go
COPY fastnode-golib ./fastnode-golib
COPY services/release.khulnasoft.com ./services/release.khulnasoft.com

WORKDIR /go/src/github.com/khulnasoft-lab/fastnode/services/release.khulnasoft.com
RUN go build .

FROM ubuntu:20.04
WORKDIR /root/
ENV RELEASE_DB_DRIVER=postgres
COPY --from=0 /go/src/github.com/khulnasoft-lab/fastnode/services/release.khulnasoft.com/release.khulnasoft.com .
ENTRYPOINT ["./release.khulnasoft.com"]

