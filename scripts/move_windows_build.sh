#!/usr/bin/env bash

# Move the staged windows build to the release repo

set -e # exit if any command fails

KHULNASOFT-LAB="${KHULNASOFT-LAB:-$HOME/khulnasoft-lab}"
RELEASE_BUILD_DIR="${KHULNASOFT-LAB}/scripts"

if [-z "$1"]; then
  echo "No staged build directory supplied"
else
  # remove previous contents
  rm -r $RELEASE_BUILD_DIR/staged-windows-build
  # mkdir if it doesn't exist
  mkdir -p $RELEASE_BUILD_DIR/staged-windows-build
  # copy over contents of staged build dir
  cp -R $1 $RELEASE_BUILD_DIR
  echo "staged release copied over"
fi