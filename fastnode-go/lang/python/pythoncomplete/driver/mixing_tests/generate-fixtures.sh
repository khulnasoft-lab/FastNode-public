#!/bin/bash

# Runs record-completions cmd on all mixing test files to generate test fixtures.

set -e

GOPATH=$(go env GOPATH)
KHULNASOFT-LAB="$GOPATH/src/github.com/khulnasoft-lab/fastnode"
COMPLETEDIR="$KHULNASOFT-LAB/fastnode-go/lang/python/pythoncomplete"
CMDDIR="$COMPLETEDIR/offline/cmds/record-completions"
INPUTSDIR="$COMPLETEDIR/driver/mixing_tests/inputs"

cd $CMDDIR

go build

for file in $INPUTSDIR/*
do
    if [[ -f $file ]]; then
        ./record-completions $file
    fi
done

rm record-completions
