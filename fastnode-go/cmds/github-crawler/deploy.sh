#!/bin/bash

for i in `seq 0 4`; do
    echo "stoping crawler-$i..."
    ssh crawler-$i.khulnasoft.com "killall github-crawler"
done

GOOS=linux go build

for i in `seq 0 4`; do
    echo "deploying to crawler-$i..."
    scp crawler.sh crawler-$i.khulnasoft.com:
    scp github-crawler crawler-$i.khulnasoft.com:
done

echo "starting master @ crawler-0..."
scp master.sh crawler-0.khulnasoft.com:
ssh crawler-0.khulnasoft.com "./master.sh"

for i in `seq 0 4`; do
    echo "starting crawler-$i..."
    ssh crawler-$i.khulnasoft.com "rm -rf outputdir/*"
    ssh crawler-$i.khulnasoft.com "./crawler.sh"
done

echo "done"
