#!/bin/bash
echo "building.."
GOOS=linux go build -o ./usernode-inspector github.com/khulnasoft-lab/fastnode/fastnode-go/user/cmds/usernode-inspector
GOOS=linux go build -o ./status-inspector github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/status-inspector
GOOS=linux go build -o ./web-snapshots github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/web-snapshots

echo "killing existing.."
ssh metrics.khulnasoft.com "killall usernode-inspector"
ssh metrics.khulnasoft.com "killall status-inspector"
ssh metrics.khulnasoft.com "killall web-snapshots"

sleep 1

echo "syncing.."
scp usernode-inspector metrics.khulnasoft.com:
scp status-inspector metrics.khulnasoft.com:
scp web-snapshots metrics.khulnasoft.com:

sleep 1

echo "starting.."

ssh metrics.khulnasoft.com 'nohup ./status-inspector -port=":4040" &> status.log &'
ssh metrics.khulnasoft.com 'nohup ./usernode-inspector -port=":3030" &> usernode.log &'
ssh metrics.khulnasoft.com 'nohup ./web-snapshots -urls=http://metrics.khulnasoft.com,http://users.khulnasoft.com &> snapshots.log &'

rm -f usernode-inspector
rm -f status-inspector
rm -f web-snapshots
