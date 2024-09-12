#!/bin/bash

set -e

VER=$(go version | cut -d" " -f3)
if [ "$VER" != "go1.15.3" ]; then
    echo "Please install Go 1.15.3"
    exit 1
fi

export KHULNASOFT-LAB="${KHULNASOFT-LAB:-$(go env GOPATH)/src/github.com/khulnasoft-lab/fastnode}"

export GOPRIVATE=github.com/khulnasoft-lab/*

# See notes at https://golang.org/cmd/cgo/#hdr-Using_cgo_with_the_go_command
export CGO_CFLAGS_ALLOW=".+"
export CGO_LDFLAGS_ALLOW=".+"

rm -rf $KHULNASOFT-LAB/osx/libfastnoded
mkdir -p $KHULNASOFT-LAB/osx/libfastnoded
cd $KHULNASOFT-LAB/osx/libfastnoded
go build \
	-buildmode=c-archive \
	-ldflags "-L $KHULNASOFT-LAB/osx/tensorflow/lib" \
	-gcflags "-I $KHULNASOFT-LAB/osx/tensorflow/include" \
	-ldflags "-X github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clientapp.gitCommit=$(git rev-parse --short HEAD)" \
	github.com/khulnasoft-lab/fastnode/fastnode-go/client/libfastnoded
