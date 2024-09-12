#!/usr/bin/env bash

set -e

export GOPATH="$PWD/gopath"
export GO111MODULE="on"
export GOPRIVATE=github.com/khulnasoft-lab/*
BUILD_DIR="$PWD/build"

KHULNASOFT-LAB=$GOPATH/src/github.com/khulnasoft-lab/fastnode

VERSION=$(cat version/version)
COMMIT=$(cat version/commit)
PREVVERSION=$(cat version/prev)

echo "Building Linux client..."
echo "VERSION=$VERSION"
echo "COMMIT=$COMMIT"
echo "PREVVERSION=$PREVVERSION"
echo

# TODO make this work without the cd. cd $KHULNASOFT-LAB isn't enough either...
cd $KHULNASOFT-LAB/linux && ./docker/build_installer.sh "$VERSION" "$BUILD_DIR"
cd $KHULNASOFT-LAB/linux && ./docker/build_update_package.sh "$VERSION" "$BUILD_DIR" "$PRIVATE_KEY" "$PREVVERSION"

echo
echo

mv "$BUILD_DIR/version-$VERSION.json" "$BUILD_DIR/version.json"
mv "$BUILD_DIR/fastnode-installer-$VERSION" "$BUILD_DIR/fastnode-installer"
mv "$BUILD_DIR/fastnode-installer-$VERSION.sh" "$BUILD_DIR/fastnode-installer.sh"
mv "$BUILD_DIR/fastnode-updater-$VERSION.sh" "$BUILD_DIR/fastnode-updater.sh"
echo "PLATFORM=linux" >> "$BUILD_DIR/META"
echo "VERSION=$VERSION" >> "$BUILD_DIR/META"
echo "COMMIT=$COMMIT" >> "$BUILD_DIR/META"
echo "SIGNATURE=" >> "$BUILD_DIR/META"
echo "build/META:"
cat $BUILD_DIR/META && echo

if [ -n "$PREVVERSION" ]; then
    mkdir -p $BUILD_DIR/deltaFrom/$PREVVERSION
    mv $BUILD_DIR/fastnode-patch-updater$PREVVERSION-$VERSION.sh $BUILD_DIR/deltaFrom/$PREVVERSION/fastnode-updater.sh
    mv "$BUILD_DIR/version-$PREVVERSION-$VERSION.json" $BUILD_DIR/deltaFrom/$PREVVERSION/version.json
    echo "DELTA_FROM[0]=$PREVVERSION" >> $BUILD_DIR/META
    echo "PLATFORM=linux" >> $BUILD_DIR/deltaFrom/$PREVVERSION/META
    echo "FROM_VERSION=$PREVVERSION" >> $BUILD_DIR/deltaFrom/$PREVVERSION/META
    echo "TO_VERSION=$VERSION" >> $BUILD_DIR/deltaFrom/$PREVVERSION/META
    echo "SIGNATURE=" >> $BUILD_DIR/deltaFrom/$PREVVERSION/META
    echo "build/deltaFrom/$PREVVERSION/META:"
    cat $BUILD_DIR/deltaFrom/$PREVVERSION/META && echo
fi
