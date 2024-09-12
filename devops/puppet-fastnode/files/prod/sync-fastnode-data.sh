#!/bin/bash

DEST=/var/fastnode/data/
BUCKET=s3://fastnode-data/var-fastnode-data/

sudo s3cmd sync $BUCKET $DEST \
    --access_key=${FASTNODE_DATA_ACCESS_KEY} \
    --secret_key=${FASTNODE_DATA_SECRET_KEY} \
    --ssl \
    --follow-symlinks \
    --recursive \
    --human-readable-sizes \
    --no-delete-removed \
    --progress

sudo chown -R $USER:$USER $DEST
