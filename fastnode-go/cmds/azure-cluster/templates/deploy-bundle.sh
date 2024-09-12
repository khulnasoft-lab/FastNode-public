#!/usr/bin/env bash

set -e

BUNDLE_FILE=/var/fastnode/upload/bundle.tar.gz
BUNDLE_PATH=/var/fastnode

tar xzvf $BUNDLE_FILE -C $BUNDLE_PATH

bash /var/fastnode/bundle/setup.sh
source /var/fastnode/bundle/env.sh

echo "export INSTANCE_ID=`cat /var/fastnode/instance-id`" >> /var/fastnode/bundle/env.sh
echo "export INSTANCE_COUNT=`cat /var/fastnode/instance-count`" >> /var/fastnode/bundle/env.sh

bash /var/fastnode/bundle/start.sh

# stop these clusters once the bundle has successfully finished
{{range $cluster := .CleanupClusters}}
echo "Stopping {{$cluster}}"
azure-cluster stop {{$cluster}}
{{end}}
