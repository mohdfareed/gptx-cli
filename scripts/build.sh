#!/usr/bin/env bash

# MARK: Init ==================================================================

# help
USAGE="usage: $0 [output=.bin]"
if [ "$#" -gt 1 ]; then echo "$USAGE" && exit 1; fi # don't bother reading

# args
APP=$(realpath ./cmd/gptx) # the app source
BIN="${1:-.bin}" # the binaries path

# setup
go build -o "$BIN/_" "$APP" # download deps
rm -rf "$BIN" && mkdir -p "$BIN" # clear bin

# MARK: Build =================================================================

# build for a plat-arch and package it
build() { # usage: build <plat> <arch> <id>
  plat="$1"; arch="$2"; id="$3";
  archive=$BIN/$(basename "$APP")-$id.zip # the archive name

  # build and package (env, build, perm, zip)
  echo "building for $plat $arch..."
  GOOS=$plat GOARCH=$arch go build -C "$BIN" "$APP"
  find "$BIN" -maxdepth 1 -type f ! -name '*.zip' -exec chmod +x {} \;
  zip -jm "$archive" "$BIN"/* -x '*.zip' > /dev/null
  echo "-> packaged: $archive"
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
echo "building for development..."
go build -C "$BIN" -tags=dev "$APP"
echo "-> dev pkg at: $BIN/$(basename "$APP")"

# update documentation
echo "updating docs..."
echo "# gptx" > "cmd/cli.md"
echo >> "cmd/cli.md"
if [ "$GOOS" = "windows" ]; then
  ./$BIN/gptx.exe -h >> "cmd/cli.md"
else
  ./$BIN/gptx -h >> "cmd/cli.md"
fi
