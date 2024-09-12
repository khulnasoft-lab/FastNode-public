#!/bin/bash

set -e

GITCOMMIT=$(git rev-parse HEAD)
export KHULNASOFT-LAB="${KHULNASOFT-LAB:-$(go env GOPATH)/src/github.com/khulnasoft-lab/fastnode}"
KHULNASOFT-LAB_LINUX="$KHULNASOFT-LAB/linux"

cd "$KHULNASOFT-LAB_LINUX"

echo "Checking for the copilot (fastnode) ..."
if [ ! -f "linux-unpacked/fastnode" ]; then
    echo "Please run 'fastnode-go/linux/build_electron.sh' to build the copilot"
    exit 1
fi

echo "Found the copilot, building fastnoded (${GITCOMMIT}) ..."

rm -f fastnoded

# at runtime our locally build fastnoded binary checks tensorflow/lib for shared libraries
# the deployed fastnoded binary is only checking . before $LD_LIBRARY_PATH (see cmds/fastnoded/main.go)
go build \
    -ldflags "-X github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clientapp.gitCommit=${GITCOMMIT}" \
    -ldflags "-r tensorflow/lib" \
    github.com/khulnasoft-lab/fastnode/fastnode-go/client/cmds/fastnoded

./fastnoded
