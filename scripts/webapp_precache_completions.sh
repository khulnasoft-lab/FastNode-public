#!/bin/bash
set -e

KHULNASOFT-LAB="${KHULNASOFT-LAB:-$GOPATH/src/github.com/khulnasoft-lab/fastnode}"
OUTDIR=$KHULNASOFT-LAB/webapp/src/assets/data/precaching

cd $KHULNASOFT-LAB/fastnode-go/cmds/completions-precacher
go generate
go build
mkdir -p $OUTDIR
./completions-precacher -out=$OUTDIR -name=completions.json