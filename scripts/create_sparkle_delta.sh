#!/usr/bin/env bash
# Usage: ./create_sparkle_delta.sh --from-dir <from_version_dir> --to-dir <to_version_dir>
# The version directories contain the Fastnode.app. They should be named the version they contain.
# The script creates a file <from_version>-<to_version>.delta in the directory from which it is run.

set -ev # exit if any command fails, verbose

KHULNASOFT-LAB="${KHULNASOFT-LAB:-$HOME/khulnasoft-lab}"

# parse args
while [[ $# > 0 ]]
do
    key="$1"
    case $key in
        --from-app)
        FROMAPP="$2"
        shift
        shift
        ;;
        --to-app)
        TOAPP="$2"
        shift
        shift
        ;;
        --delta)
        DELTA="$2"
        shift
        shift
        ;;
        *)
        # unknown option
        echo "Unknown option:" $key
        exit 1
        ;;
    esac
done

if [[ -z "$FROMAPP" ]]; then
    echo "Location of \"from\" Fastnode.app is not set, set it using --from-app. exiting."
    exit 1
fi

if [[ -z "$TOAPP" ]]; then
    echo "Location of \"to\" Fastnode.app is not set, set it using --to-app. exiting."
    exit 1
fi

if [[ -z "$DELTA" ]]; then
    echo "Name of delta archive to output is not set, set it using --delta. exiting."
    exit 1
fi

$KHULNASOFT-LAB/osx/BinaryDelta create --verbose $FROMAPP $TOAPP $DELTA
