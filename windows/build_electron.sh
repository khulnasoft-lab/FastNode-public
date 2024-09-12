#!/bin/bash

set -e
KHULNASOFT-LAB=${KHULNASOFT-LAB:-$(go env GOPATH)/src/github.com/khulnasoft-lab/fastnode}
ENVIRONMENT="${REACT_APP_ENV:-development}"

cd $KHULNASOFT-LAB/sidebar
rm -rf dist
npm install
REACT_APP_ENV=$ENVIRONMENT npm run pack:win
rm -rf $KHULNASOFT-LAB/windows/win-unpacked
cp -r $KHULNASOFT-LAB/sidebar/dist/win-unpacked $KHULNASOFT-LAB/windows/
rm -rf $KHULNASOFT-LAB/sidebar/dist
