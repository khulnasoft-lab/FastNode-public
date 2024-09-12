#!/usr/bin/env bash
set -e
export KHULNASOFT-LAB="${KHULNASOFT-LAB:-$(go env GOPATH)/src/github.com/khulnasoft-lab/fastnode}"
export ENVIRONMENT=development

cd $KHULNASOFT-LAB

if [[ $CONFIGURATION == "Release" ]]; then
    ENVIRONMENT=production
fi

if [[ $1 == "force" || $CONFIGURATION == "Release" ]]; then
    rm -rf dist
    cd $KHULNASOFT-LAB/sidebar
    npm install
    echo "ENV: $ENVIRONMENT"
    REACT_APP_ENV=$ENVIRONMENT yarn run pack
    exit 0
fi

echo "Checking for electon/Fastnode.app..."
if [ ! -d "sidebar/dist/mac/Fastnode.app" ]; then
    echo "... not found. Please build the sidebar application."
    exit 1
fi

echo "... found!"
