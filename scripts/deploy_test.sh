#!/bin/bash

BUILD=/tmp/build
mkdir -p $BUILD

if [ "$#" -ne 1 ]; then
    echo "1 argument required (e.g test-N.khulnasoft.com)"
    exit 1
fi

if [ -z "$SKIP_LOCAL_CODE_WORKER_BUILD" ]; then
    echo "== building local-code-worker..."
    GOOS=linux go build -o $BUILD/local-code-worker github.com/khulnasoft-lab/fastnode/fastnode-go/localcode/cmds/local-code-worker 2> >(sed $'s,.*,\e[31m&\e[m,'>&2) || exit 1
else
    echo "skipping local-code-worker build..."
fi

if [ -z "$SKIP_USER_NODE_BUILD" ]; then
    echo "== building user-node..."
    GOOS=linux go build -o $BUILD/user-node github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/user-node 2> >(sed $'s,.*,\e[31m&\e[m,'>&2) || exit 1
else
    echo "skipping user-node build..."
fi

if [ -z "$SKIP_USER_MUX_BUILD" ]; then
    echo "== building user-mux..."
    GOOS=linux go build -o $BUILD/user-mux github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/user-mux 2> >(sed $'s,.*,\e[31m&\e[m,'>&2) || exit 1
else
    echo "skipping user-mux build..."
fi

echo "== killing existing processes..."
if [ -z "$SKIP_LOCAL_CODE_WORKER_BUILD" ]; then
    ssh $1 "killall local-code-worker"
fi

if [ -z "$SKIP_USER_NODE_BUILD" ]; then
    ssh $1 "killall user-node"
fi

if [ -z "$SKIP_USER_MUX_BUILD" ]; then
    ssh $1 "killall user-mux"
fi

sleep 5

echo "== syncing binaries..."
rsync -r $BUILD/* $1:/deploy/ || exit 1

if [ -z "$SKIP_LOCAL_CODE_WORKER_BUILD" ]; then
    echo "== reloading local-code-worker..."
    ssh $1 "nohup /deploy/local-code-worker &> /var/fastnode/log/local-code-worker.log < /dev/null &" || exit 1
fi

if [ -z "$SKIP_USER_MUX_BUILD" ]; then
    echo "== reloading user-mux..."
    ssh $1 "nohup /deploy/user-mux &> /var/fastnode/log/user-mux.log < /dev/null &" || exit 1
fi

if [ -z "$SKIP_USER_NODE_BUILD" ]; then
    echo "== reloading user-node..."
    ssh $1 "nohup /deploy/user-node $USERNODE_ARGS &> /var/fastnode/log/user-node.log < /dev/null &" || exit 1

    echo -n "== waiting for user-node to load"
    if [ -z "$TAIL_LOG" ]; then
        until `curl --output /dev/null --silent --head --fail http://$1:9081/ready`; do
            printf '.'
            sleep 5
        done
        printf '\n'
    else
        ssh $1 "tail -f /var/fastnode/log/user-node.log" || exit 1
    fi
fi

echo "== done!"
