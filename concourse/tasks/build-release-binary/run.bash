#!/usr/bin/env bash

set -e

export GOPATH=$PWD/gopath
export BUILD_DIR=$PWD
export GO111MODULE="on"

KHULNASOFT-LAB=$GOPATH/src/github.com/khulnasoft-lab/fastnode
cd $KHULNASOFT-LAB

go build -o $BUILD_DIR/release_bin/release github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/release
