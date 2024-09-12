#!/usr/bin/env bash

# Script for running the various make commands to generate bindata

# NOTE: part of the reason why this is here and not in slackbuildbot is because make seems to be
# unhappy about being run from subprocess.Popen and doesn't get run in the correct working
# directory

set -e

# set KHULNASOFT-LAB to $HOME/khulnasoft-lab if it's not already specified
KHULNASOFT-LAB="${KHULNASOFT-LAB:-$HOME/khulnasoft-lab}"
echo "make_bindata.sh: KHULNASOFT-LAB set to ${KHULNASOFT-LAB}"
cd $KHULNASOFT-LAB
