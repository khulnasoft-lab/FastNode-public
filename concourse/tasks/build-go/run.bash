#!/usr/bin/env bash

set -e

export GOPATH="$PWD/gopath"
BUILD_DIR="$PWD/build"
KHULNASOFT-LAB=$GOPATH/src/github.com/khulnasoft-lab/fastnode


# copied from prepare-release: TODO(naman) factor this out
tag=$(cd $KHULNASOFT-LAB && git describe --tags --exact-match)
if [[ $tag =~ ^v[0-9]{8}\.[0-9]+$ ]]; then
    true
else
    echo "Error preparing release: invalid tag format ($tag). Aborting."
    exit 1
fi
VERSION=$tag
COMMIT=$(cd $KHULNASOFT-LAB && git rev-parse HEAD)

echo "Building backend..."
echo "VERSION=$VERSION"
echo "COMMIT=$COMMIT"

CGO_LDFLAGS_ALLOW=".*" go build -o "build/$2-$VERSION" $1

echo "VERSION=$VERSION" >> build/META
echo "COMMIT=$COMMIT" >> build/META
