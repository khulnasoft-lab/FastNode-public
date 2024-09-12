#!/usr/bin/env bash

set -e

export GOPATH="$PWD/gopath"
export GO111MODULE="on"
BUILD_DIR="$PWD/build"
KHULNASOFT-LAB=$GOPATH/src/github.com/khulnasoft-lab/fastnode

cd $KHULNASOFT-LAB

# copied from prepare-release: TODO(naman) factor this out
tag=$(git describe --tags --exact-match)
if [[ $tag =~ ^v[0-9]{8}\.[0-9]+$ ]]; then
    true
else
    echo "Error preparing release: invalid tag format ($tag). Aborting."
    exit 1
fi
VERSION=$tag
COMMIT=$(git rev-parse HEAD)


echo "Building backend..."
echo "VERSION=$VERSION"
echo "COMMIT=$COMMIT"


CGO_LDFLAGS_ALLOW=".*" go build -o $BUILD_DIR/user-node github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/user-node
CGO_LDFLAGS_ALLOW=".*" go build -o $BUILD_DIR/user-mux github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/user-mux
CGO_LDFLAGS_ALLOW=".*" go build -o $BUILD_DIR/release github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/release

pushd $KHULNASOFT-LAB/fastnode-server
make fastnode-server.tgz DIR=fastnode-server TOKEN=placeholder
popd
cp $KHULNASOFT-LAB/fastnode-server/fastnode-server.tgz $BUILD_DIR

echo "VERSION=$VERSION" >> $BUILD_DIR/META
echo "COMMIT=$COMMIT" >> $BUILD_DIR/META
