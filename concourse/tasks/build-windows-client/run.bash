# Cygwin Bash.exe

set -e

export GOPATH=$(cygpath -w "$PWD/gopath")
export GO111MODULE="on"
export GOPRIVATE=github.com/khulnasoft-lab/*
KHULNASOFT-LAB=$PWD/gopath/src/github.com/khulnasoft-lab/fastnode/

COMMIT=$(cat version/commit)
VERSION=$(cat version/version)
PREVVERSION=$(cat version/prev)

echo "Building Windows client..."
echo "VERSION=$VERSION"
echo "COMMIT=$COMMIT"
echo "PREVVERSION=$PREVVERSION"
echo

REACT_APP_ENV=production $KHULNASOFT-LAB/windows/build_electron.sh
WINDOWS_BUILD_VERSION=$VERSION WINDOWS_PATCH_BASE=$PREVVERSION make -f $KHULNASOFT-LAB/Makefile -C $KHULNASOFT-LAB FastnodeSetup.exe FastnodeUpdateInfo.xml FastnodePatchUpdateInfo.xml

echo && echo

OUTDIR=$KHULNASOFT-LAB/windows/installer/builds/$VERSION

mv $OUTDIR/FastnodeSetup$VERSION.exe build/FastnodeSetup.exe
mv $OUTDIR/FastnodeUpdater$VERSION.exe build/FastnodeUpdater.exe
mv $OUTDIR/FastnodeUpdateInfo.xml build/FastnodeUpdateInfo.xml
echo "PLATFORM=windows" >> build/META
echo "VERSION=$VERSION" >> build/META
echo "COMMIT=$COMMIT" >> build/META
echo "SIGNATURE=" >> build/META
echo "build/META:"
cat build/META && echo

if [ -n "$PREVVERSION" ]; then
    mkdir -p build/deltaFrom/$PREVVERSION
    mv $OUTDIR/FastnodePatchUpdater$PREVVERSION-$VERSION.exe build/deltaFrom/$PREVVERSION/FastnodeDeltaUpdater.exe
    mv $OUTDIR/FastnodePatchUpdateInfo$PREVVERSION.xml build/deltaFrom/$PREVVERSION/FastnodeDeltaUpdateInfo.xml
    echo "DELTA_FROM[0]=$PREVVERSION" >> build/META
    echo "PLATFORM=windows" >> build/deltaFrom/$PREVVERSION/META
    echo "FROM_VERSION=$PREVVERSION" >> build/deltaFrom/$PREVVERSION/META
    echo "TO_VERSION=$VERSION" >> build/deltaFrom/$PREVVERSION/META
    echo "SIGNATURE=" >> build/deltaFrom/$PREVVERSION/META
    echo "build/deltaFrom/$PREVVERSION/META:"
    cat build/deltaFrom/$PREVVERSION/META
fi
