#!/bin/bash

SELF_DIR=$(dirname $BASH_SOURCE)

# If under linux then build natively
if go version | grep -q linux; then
	go build github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/user-node
	exit
fi

echo "Creating khulnasoft-lab.tar.gz..."
rm -f /tmp/khulnasoft-lab.tar.gz
tar czf /tmp/khulnasoft-lab.tar.gz \
	-C $GOPATH/src \
	github.com/khulnasoft-lab/fastnode/fastnode-go \
	github.com/khulnasoft-lab/fastnode/fastnode-golib \
	github.com/khulnasoft-lab/fastnode/vendor

echo "Creating container..."
CONTAINER=$(docker create -e GO15VENDOREXPERIMENT=1 golang:1.6 go build -o /user-node github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/user-node)
echo "Created $CONTAINER"

echo "Copying khulnasoft-lab source into container..."
cat /tmp/khulnasoft-lab.tar.gz | docker cp - $CONTAINER:/go/src/

echo "Starting container..."
docker start --attach $CONTAINER || exit 1

echo "Retrieving binary..."
docker cp $CONTAINER:/user-node user-node || exit 1

echo "Cleaning up..."
docker rm -f $CONTAINER || exit 1
rm -f /tmp/khulnasoft-lab.tar.gz
