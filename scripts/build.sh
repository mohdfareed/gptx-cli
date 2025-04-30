#!/bin/sh

# help message
USAGE="usage: $0 [path=.bin]"
if [ "$#" -gt 1 ]; then echo "$USAGE" && exit 1; fi

# arguments
APP=./chatgpt # the app source code
BIN="${1:-.bin}" # the binaries path
exec=$(basename "$APP") # the executable
out="$BIN/$exec" # the output path

# build for a plat-arch and package it
build() { # usage: build <plat> <arch> <id>
  plat="$1"; arch="$2"; id="$3"

  # build the executable
  echo "building for $id..."
  GOOS=$plat GOARCH=$arch go build -o "$out" "$APP"

  # package into an archive
  archive="$out-$id.zip"
  zip -j "$archive" "$out" > /dev/null
}

# linux
build linux arm64 "linux-arm"
build linux amd64 "linux-x64"
# macos
build darwin arm64 "macos-arm"
build darwin amd64 "macos-x64"
# windows
build windows arm64 "win-arm"
build windows amd64 "win-x64"

# cleanup
rm "$out"
