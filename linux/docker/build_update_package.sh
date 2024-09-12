#!/usr/bin/env bash
set -e

command -v makeself >/dev/null 2>&1 || { echo >&2 "makeself is required to build the fastnode-install package.  Aborting."; exit 1; }
command -v chrpath >/dev/null 2>&1 || { echo >&2 "chrpath is required to build the fastnode-install package.  Aborting."; exit 1; }
command -v git >/dev/null 2>&1 || { echo >&2 "git is required to build the fastnode-install package.  Aborting."; exit 1; }
command -v sha256sum >/dev/null 2>&1 || { echo >&2 "sha256sum is required to build the fastnode-install package.  Aborting."; exit 1; }

VERSION="$1"
[ -n "$VERSION" ] || { echo >&2 "version parameter not passed"; exit 1; }

BUILD_DIR="$2"
[ -n "$BUILD_DIR" ] || { echo >&2 "build dir parameter not passed"; exit 1; }

PRIVATE_KEY="$3"
[ -n "$PRIVATE_KEY" ] || { echo >&2 "private key parameter not passed"; exit 1; }

PREVVERSION="$4"
[ -n "$PREVVERSION" ] || { echo >&2 "previous version parameter not passed"; exit 1; }

PAR=$(dirname "$(readlink -f "$0")")
ROOT="$PAR/fastnode-v$VERSION"
GITCOMMIT=$(git rev-parse HEAD)
CURR_DIR="$ROOT/fastnode-v$VERSION"
PATCH_DIR="$ROOT/patch"

mkdir -p "$ROOT"
mkdir -p "$CURR_DIR" "$PATCH_DIR"

echo "Building fastnoded $VERSION (git commit $GITCOMMIT)..."
CGO_LDFLAGS_ALLOW=".*" go build \
    -ldflags "-X github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clientapp.gitCommit=${GITCOMMIT} -X github.com/khulnasoft-lab/fastnode/fastnode-go/client/platform/version.version=${VERSION} -r tensorflow/lib" \
    -o $CURR_DIR/fastnoded \
    github.com/khulnasoft-lab/fastnode/fastnode-go/client/cmds/fastnoded

echo "Copying tensorflow libraries to $CURR_DIR/lib..."
cp -Ra ./tensorflow/lib "$CURR_DIR"

echo "Changing rpath of fastnoded binary..."
chrpath -r '$ORIGIN/lib:$ORIGIN' "$CURR_DIR/fastnoded"

echo "Building electron application..."
export REACT_APP_END=production
bash ./build_electron.sh
cp -Ra linux-unpacked "$CURR_DIR/"

echo "Building fastnode-lsp binary..."
go build \
    -o "$CURR_DIR/fastnode-lsp" \
    github.com/khulnasoft-lab/fastnode/fastnode-go/lsp/cmds/fastnode-lsp

echo "Building fastnode updater binary..."
go build \
    -o "$CURR_DIR/fastnode-update" \
    -ldflags "-X main.version=${VERSION}" \
    github.com/khulnasoft-lab/fastnode/linux/cmds/fastnode-installer

echo "Creating updater executable..."
TARGET_FILE="$BUILD_DIR/fastnode-updater-$VERSION.sh"
makeself --notemp --nox11 "$CURR_DIR" "$TARGET_FILE" "Fastnode updater, version $VERSION"

UPDATER_CHECKSUM="$(sha256sum "$TARGET_FILE" | cut -d " " -f 1)"

go build -o fastnode-update-signer ./cmds/fastnode-update-signer/
UPDATER_SIGNATURE="$(./fastnode-update-signer "$TARGET_FILE" "$PRIVATE_KEY")"

cat - <<EOF > "$BUILD_DIR/version-$VERSION.json"
{
    "version":"$VERSION",
    "updater_url":"https://linux.khulnasoft.com/linux/$VERSION/fastnode-updater.sh",
    "sha256":"$UPDATER_CHECKSUM",
    "signature": "$UPDATER_SIGNATURE"
}
EOF

# use -s -w to strip debugging information, which reduces the binary size.
go build -ldflags="-s -w" ./cmds/tar/

echo "Creating current version archive..."
./tar archive "$CURR_DIR" > "$ROOT/$VERSION.tar" 

echo "Getting previous version archive..."
LINUX_BUCKET='s3://fastnode-downloads/linux'
KEY="$PREVVERSION/fastnode-updater.sh"
aws s3 ls "$LINUX_BUCKET/$KEY" || not_exist=true
if [ $not_exist ]
then
    echo "$LINUX_BUCKET/$KEY does not exist, skipping creating delta update..."
    exit 0
fi
aws s3 cp "$LINUX_BUCKET/$KEY" "$ROOT/fastnode-updater.sh"

echo "Extracting previous version..."
chmod u+x "$ROOT/fastnode-updater.sh"
$ROOT/fastnode-updater.sh --noexec --target "$ROOT/$PREVVERSION"

echo "Creating tar of previous version..."
# these flags are necessary to ensure the tar is the same when we create the patch and when
# we apply it
./tar archive "$ROOT/$PREVVERSION" > "$ROOT/$PREVVERSION.tar"

PREVVERSION_CHECKSUM="$(sha256sum "$ROOT/$PREVVERSION.tar" | cut -d " " -f 1)"
echo "Checksum of previous version tar: $PREVVERSION_CHECKSUM"

# get the version of go-bsdiff specified in github.com/khulnasoft-lab/fastnode/linux/go.mod
go get -v github.com/khulnasoft-lab/go-bsdiff/v2
# use -s -w to strip debugging information, which reduces the binary size.
go build -ldflags="-s -w" github.com/khulnasoft-lab/go-bsdiff/v2/cmd/bsdiff
go build -ldflags="-s -w" github.com/khulnasoft-lab/go-bsdiff/v2/cmd/bspatch

echo "Creating patch file..."
./bsdiff "$ROOT/$PREVVERSION.tar" "$ROOT/$VERSION.tar" "$PATCH_DIR/$PREVVERSION-$VERSION.patch"

echo "Verifying patch file applies correctly..."
./bspatch "$ROOT/$PREVVERSION.tar" "$ROOT/$VERSION.new.tar" "$PATCH_DIR/$PREVVERSION-$VERSION.patch"
DIFF="$(diff $ROOT/$VERSION.tar $ROOT/$VERSION.new.tar)"
if [ "$DIFF" ]
then
    echo "The patched update does not match the new version..."
    exit 0
fi

echo "Copying apply patch script to $PATCH_DIR..."
cp -p apply_patch.sh "$PATCH_DIR"

echo "Copying bspatch and tar binaries to $PATCH_DIR..."
cp -p bspatch "$PATCH_DIR"
cp -p tar "$PATCH_DIR"

echo "Creating patch updater executable..."
PD="$(basename $PATCH_DIR)"
PATCH_TARGET_FILE="$BUILD_DIR/fastnode-patch-updater$PREVVERSION-$VERSION.sh"
makeself --notemp --nox11 "$PATCH_DIR" "$PATCH_TARGET_FILE" "Fastnode patch updater, version $PREVVERSION-$VERSION" ./apply_patch.sh $PREVVERSION-$VERSION.patch "$PREVVERSION" "$VERSION" "$PD"

PATCH_UPDATER_CHECKSUM="$(sha256sum "$PATCH_TARGET_FILE" | cut -d " " -f 1)"

PATCH_UPDATER_SIGNATURE="$(./fastnode-update-signer "$PATCH_TARGET_FILE" "$PRIVATE_KEY")"

cat - <<EOF > "$BUILD_DIR/version-$PREVVERSION-$VERSION.json"
{
    "version":"$VERSION",
    "updater_url":"https://linux.khulnasoft.com/linux/$VERSION/$PREVVERSION/fastnode-patch-updater.sh",
    "sha256":"$PATCH_UPDATER_CHECKSUM",
    "signature": "$PATCH_UPDATER_SIGNATURE"
}
EOF
