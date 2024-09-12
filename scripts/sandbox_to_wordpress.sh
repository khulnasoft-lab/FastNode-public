#!/bin/bash
set -e

KHULNASOFT-LAB="${KHULNASOFT-LAB:-$GOPATH/src/github.com/khulnasoft-lab/fastnode}"

cd $KHULNASOFT-LAB/web/sandbox
# remove previous assets
rm -rf ./dist

# build assets
npm run make-production-all

cd dist

# reset known_hosts entry
ssh-keygen -R [XXXXXXX]:2222

# sync remote directory with current
# requires a LFTP_PASSWORD environment variable
lftp --env-password \
      -p 2222 \
      sftp://XXXXXXX@XXXXXXX.XXXXXXX.XXXXXXX.com \
      -e "set sftp:auto-confirm yes; cd wp-content/fastnode-sandbox; mirror -R --delete-first; ls; bye"
