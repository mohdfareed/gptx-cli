#!/bin/sh

# MARK: Arguments =============================================================

# help message
USAGE="usage: $0 [path=.bin]"
if [ "$#" -gt 1 ]; then echo "$USAGE" && exit 1; fi

APP=./gptx # the app source code
BIN="${1:-.bin}" # the binaries path

exec=$(basename "$APP") # the executable
out="$BIN/$exec" # the output path

# MARK: Build =================================================================

# build for a plat-arch and package it
build() { # usage: build <plat> <arch> <id>
  plat="$1"; arch="$2"; id="$3"
  output="$out"

  # add .exe for windows
  if [ "$plat" = "windows" ]; then
    output="$output.exe"
  fi

  # build the executable
  echo "building for $id..."
  GOOS=$plat GOARCH=$arch go build -o "$output" "$APP"

  # package into an archive
  archive="$out-$id.zip"
  zip -j "$archive" "$output" > /dev/null
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
if [ "$(go env GOOS)" = "windows" ]; then
  go build -tags=debug -o "$out.exe" "$APP"
  rm "$out"
else # macos | linux
  go build -tags=debug -o "$out" "$APP"
  rm "$out.exe"
fi

echo "builds created in $BIN"
