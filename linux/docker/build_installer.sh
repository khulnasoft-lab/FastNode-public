#!/usr/bin/env bash
# This script is supposed to be run in an Ubuntu docker container

set -e

VERSION="$1"
[ -n "$VERSION" ] || { echo >&2 "version not passed. Aborting."; exit 1; }

BUILD_DIR="$2"
[ -n "$BUILD_DIR" ] || { echo >&2 "build dir not passed. Aborting."; exit 1; }

PAR=$(dirname "$(readlink -f "$0")")
cd "$PAR"

echo "Building fastnode-installer binary..."
go build \
    -o "$BUILD_DIR/fastnode-installer-$VERSION" \
    -ldflags "-X main.version=${VERSION}" \
    github.com/khulnasoft-lab/fastnode/linux/cmds/fastnode-installer

echo "Creating fastnode-installer.sh wrapper script..."
cp "./fastnode-installer.sh" "$BUILD_DIR/fastnode-installer-$VERSION.sh"
