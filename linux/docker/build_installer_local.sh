#!/bin/bash
# This script is to build a fastnode-installer binary for local testing
# It's built in a docker container for a wider binary compatibility

set -e
DIR=$PWD
BUILD_DIR="$DIR/build"

rm -rf "$BUILD_DIR"
mkdir "$BUILD_DIR"
cat >"$BUILD_DIR/run.sh" <<EOF
#!/bin/bash
set -e
echo "Building fastnode in docker..."
export GOPATH="/src"
cd /src
go build -o "/build/fastnode-installer" github.com/khulnasoft-lab/fastnode/linux/cmds/fastnode-installer
echo Successfully built fastnode-installer.
echo
EOF

chmod +x "$BUILD_DIR/run.sh"

docker run --rm -t \
  --mount "src=$BUILD_DIR,target=/build,type=bind" \
  --mount "src=${GOPATH:-"$HOME/go"},target=/src,type=bind" \
  golang:1.14.3 /build/run.sh
