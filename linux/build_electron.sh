#!/bin/bash

set -e
export KHULNASOFT-LAB="${KHULNASOFT-LAB:-$(go env GOPATH)/src/github.com/khulnasoft-lab/fastnode}"
export ENVIRONMENT="${REACT_APP_ENV:-development}"

cd $KHULNASOFT-LAB/sidebar
rm -rf dist
npm install
REACT_APP_ENV=$ENVIRONMENT npm run pack:linux
rm -rf $KHULNASOFT-LAB/linux/linux-unpacked
cp -r $KHULNASOFT-LAB/sidebar/dist/linux-unpacked $KHULNASOFT-LAB/linux/
rm -rf $KHULNASOFT-LAB/sidebar/dist
