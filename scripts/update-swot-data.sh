#!/bin/bash
set -e
command -v "kotlinc" >/dev/null 2>&1 || { echo >&2 "kotlinc is required."; exit 1; }
command -v "java" >/dev/null 2>&1 || { echo >&2 "java is required."; exit 1; }

DIR="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd -P)"
cd "$DIR"

echo "Cloning fastnode-swot..."
rm -rf "fastnode-swot"
git clone -q "https://github.com/khulnasoft-lab/swot.git" fastnode-swot

echo "Compiling fastnode-swot..."
cd fastnode-swot
kotlinc -nowarn src/swot/Compiler.kt src/swot/Swot.kt -include-runtime -d swot.jar
java -jar ./swot.jar

test -f "out/artifacts/swot.txt" || { echo >&2 "swot.txt wasn't found."; exit 1; }
cp "out/artifacts/swot.txt" "$DIR/../fastnode-go/community/student/cmds/updatedomains/swot.txt"

echo "Compiling updatedomains"
cd "$DIR/../fastnode-go/community/student/cmds/updatedomains"
go build -o "updatedomains" .

echo "Updating swot.txt on AWS..."
./updatedomains "$DIR/../fastnode-go/community/student/cmds/updatedomains/swot.txt"

cd "$DIR"
rm "$DIR/../fastnode-go/community/student/cmds/updatedomains/"{updatedomains,swot.txt}
rm -rf "$DIR/fastnode-swot"