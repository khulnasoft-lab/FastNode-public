#!/usr/bin/env bash

# JSON_STATS=$PWD/fastnode-completions-performance.json
# go build ./fastnode-go/lang/python/pythoncomplete/offline/cmds/performancetest
# ./performancetest --json "$JSON_STATS" >/dev/null 2>&1
# gzip "$JSON_STATS"
# aws s3 cp "${JSON_STATS}.gz" "s3://fastnode-offline-metrics/model-performance/$(date -u +%Y-%m-%d-%H%M%S).json.gz"

# RECALL_STATS=$PWD/fastnode-completions-recall.json
# go build ./fastnode-go/lang/python/pythoncomplete/offline/cmds/recalltest
# ./recalltest fastnode-go/lang/python/pythoncomplete/recalltest/samples_3347.json 1000 --json "$RECALL_STATS" >/dev/null 2>&1
# cp $RECALL_STATS ${TRAVIS_BUILD_DIR}/recalls_$(date -u +%Y_%m_%d).json # We copy the recall results to build dir to try to avoid recomputing them during the go unit test
# gzip "$RECALL_STATS"

# aws s3 cp "${RECALL_STATS}.gz" "s3://fastnode-offline-metrics/model-recall/$(date -u +%Y-%m-%d-%H%M%S).json.gz"

slow_pkgs=$(find ./fastnode-go ./fastnode-golib ./local-pipelines -type f -and -iname \*.go -and -exec grep -q '+build slow' {} \; -and -exec bash -c 'echo $(dirname {})/' \;)
scripts/limitci_linux.sh '\.go$' "go test -v -tags slow -timeout 60m $slow_pkgs"

# this is failing, @juan to look (see #11599)
# sudo apt-get install python3-venv
# local-pipelines/lexical/train/scripts/ci_test_searchconfiggen.sh
