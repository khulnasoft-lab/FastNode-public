#!/usr/bin/env bash

KHULNASOFT-LAB="$HOME/khulnasoft-lab"

cd $KHULNASOFT-LAB/osx
rm -rf build/
mkdir -p build
xcodebuild -scheme Fastnode -configuration Debug -derivedDataPath build
cd build/Build/Products/Debug/
zip -r Fastnode.zip Fastnode.app
python -m SimpleHTTPServer
