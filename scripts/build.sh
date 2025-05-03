#!/usr/bin/env bash
shopt -s extglob

# MARK: Arguments =============================================================

# help message
USAGE="usage: $0 [path=.bin]"
if [ "$#" -gt 1 ]; then echo "$USAGE" && exit 1; fi # don't bother reading

APP=$(realpath ./gptx) # the app source code
BIN="${1:-.bin}" # the binaries path
mkdir -p "$BIN" # create the binaries path

# MARK: Build =================================================================

# build for a plat-arch and package it
build() { # usage: build <plat> <arch> <id>
  plat="$1"; arch="$2"; id="$3";
  archive=$BIN/$(basename "$APP")-$id.zip # the archive name

  # build and package
  echo "building for $plat $arch..."
  GOOS=$plat GOARCH=$arch go build -C "$BIN" "$APP"
  zip -jm "$archive" "$BIN"/!(*.zip) > /dev/null # zip the binaries
  echo "  packaged: $archive"
}

# MARK: Targets ===============================================================

# linux
build linux arm64 "linux-arm"
build linux amd64 "linux-x64"

# macos
build darwin arm64 "macos-arm"
build darwin amd64 "macos-x64"

# windows
build windows arm64 "win-arm"
build windows amd64 "win-x64"

# development
echo "building for debug..."
go build -C "$BIN" -tags=debug "$APP"
echo "  run $BIN/gptx --help for usage"

shopt -u extglob
