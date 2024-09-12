#!/usr/bin/env bash
set -e
go build github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/fastnodelocal/cmds/build-datadeps-filemap
./build-datadeps-filemap -verify
rm -f build-datadeps-filemap
