#!/usr/bin/env bash
set -e

BUILD_DIR="$PWD/build"
KHULNASOFT-LAB=khulnasoft-lab

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

echo "Uploading Puppet build..."
aws s3 cp $BUILD_DIR/puppet.tar.gz s3://fastnode-deploys/$VERSION/puppet.tar.gz
