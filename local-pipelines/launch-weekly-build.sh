#!/bin/bash
KHULNASOFT-LAB=$HOME/go/src/github.com/khulnasoft-lab/fastnode
#KHULNASOFT-LAB=/vagrant_home/go/src/github.com/khulnasoft-lab/fastnode
LOGDIR=$KHULNASOFT-LAB/local-pipelines/log

rm -rf $LOGDIR
mkdir -p $LOGDIR

echo "installing rundb"
go install github.com/khulnasoft-lab/fastnode/fastnode-golib/pipeline/cmds/rundb

ROOT=`rundb create s3://fastnode-data/run-db weekly-expr-model-build`
echo "using rundb root of $ROOT"

cd $KHULNASOFT-LAB/local-pipelines/python-ggnn-expr-completion

echo "starting expr-shard1, see logs at log/expr-shard1.log..."
./scripts/launch-train-shard.sh $ROOT ./packagelist-cluster1.txt expr-shard1 &> $LOGDIR/expr-shard1.log &

echo "starting expr-shard1, see logs at log/expr-shard2.log..."
./scripts/launch-train-shard.sh $ROOT ./packagelist-cluster2.txt expr-shard2 &> $LOGDIR/expr-shard2.log &

echo "waiting for expr-shard1 to complete..."
rundb wait-artifact $ROOT expr-shard1/DONE

echo "waiting for expr-shard2 to complete..."
rundb wait-artifact $ROOT expr-shard2/DONE

echo "installing buildshards"
go install github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonexpr/cmds/buildshards

echo "building shards.json"
buildshards --rundbRoot=$ROOT --output=shards.json &> $LOGDIR/buildshards.log
rundb add-artifact $ROOT shards.json shards.json
rm -rf shards.json

cd $KHULNASOFT-LAB/local-pipelines/python-call-filtering
./scripts/launch-train.sh $ROOT &> $LOGDIR/call-filtering.log &

rundb wait-artifact $ROOT call-filtering/DONE
