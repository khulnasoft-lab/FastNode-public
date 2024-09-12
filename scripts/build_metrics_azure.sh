#!/bin/bash
echo "building.."
GOOS=linux go build -o ./usernode-inspector github.com/khulnasoft-lab/fastnode/fastnode-go/user/cmds/usernode-inspector
GOOS=linux go build -o ./status-inspector github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/status-inspector
GOOS=linux go build -o ./web-snapshots github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/web-snapshots
GOOS=linux go build -o ./systems-monitor github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/systems-monitor

echo "killing existing.."
ssh metrics-azure.khulnasoft.com "killall usernode-inspector"
ssh metrics-azure.khulnasoft.com "killall status-inspector"
ssh metrics-azure.khulnasoft.com "killall web-snapshots"
ssh metrics-azure.khulnasoft.com "killall systems-monitor"

sleep 1

echo "syncing.."
scp usernode-inspector metrics-azure.khulnasoft.com:
scp status-inspector metrics-azure.khulnasoft.com:
scp web-snapshots metrics-azure.khulnasoft.com:
scp systems-monitor metrics-azure.khulnasoft.com:

sleep 1

echo "starting.."

ssh metrics-azure.khulnasoft.com 'nohup ./status-inspector -port=":4040" &> status.log &'
ssh metrics-azure.khulnasoft.com 'nohup ./usernode-inspector -port=":3030" &> usernode.log &'
ssh metrics-azure.khulnasoft.com 'nohup ./web-snapshots -urls=http://metrics-azure.khulnasoft.com,http://users-azure.khulnasoft.com &> snapshots.log &'
ssh metrics-azure.khulnasoft.com 'nohup ./systems-monitor &> monitor.log &'

rm -f usernode-inspector
rm -f status-inspector
rm -f web-snapshots
rm -f systems-monitor
