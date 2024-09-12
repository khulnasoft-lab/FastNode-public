#!/bin/bash

GITCOMMIT=$(git rev-parse HEAD)
KHULNASOFT-LAB=${KHULNASOFT-LAB:-$(go env GOPATH)/src/github.com/khulnasoft-lab/fastnode}
KHULNASOFT-LAB_WINDOWS=$KHULNASOFT-LAB/windows

cd $KHULNASOFT-LAB_WINDOWS

echo "Checking for the copilot (Fastnode.exe) ..."
if [ ! -f "win-unpacked/Fastnode.exe" ]; then
    echo "Please run 'fastnode-go/windows/build_electron.sh' to build the copilot"
    exit 1
fi


echo "Found the copilot, building fastnoded.exe (${GITCOMMIT}) ..."

# This copy of tensorflow is ignored by the .gitignore file in this directory.
# Its needed so that fastnoded.exe can launch successfully because tensorflow is dynamically linked
# (currently directory is part of the DLL search path). Always update it just incase.
cp tensorflow/lib/tensorflow.dll ./

rm -f fastnoded.exe

go build \
    -buildmode=exe \
    -ldflags "-X github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clientapp.gitCommit=${GITCOMMIT} -X github.com/khulnasoft-lab/fastnode/fastnode-go/client/sidebar.copilotDevDir=${KHULNASOFT-LAB_WINDOWS}" \
    github.com/khulnasoft-lab/fastnode/fastnode-go/client/cmds/fastnoded

./fastnoded.exe
