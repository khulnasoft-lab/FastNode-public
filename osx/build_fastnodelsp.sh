#!/bin/bash

set -e
export KHULNASOFT-LAB="${KHULNASOFT-LAB:-$GOPATH/src/github.com/khulnasoft-lab/fastnode}"

rm -rf $KHULNASOFT-LAB/osx/fastnodelsp
mkdir -p $KHULNASOFT-LAB/osx/fastnodelsp
cd $KHULNASOFT-LAB/osx/fastnodelsp
go build \
	github.com/khulnasoft-lab/fastnode/fastnode-go/lsp/cmds/fastnode-lsp
