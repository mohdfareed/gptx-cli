#!/bin/sh

# USAGE:
#    scripts/build.sh [output]
# ARGS:
#    output - the output path (default: .bin/chat)

app=./chat # the app source code
exec="${1:-.bin/chat}" # the built executable path

# the supported platforms and architectures
platforms="darwin linux windows"
architectures="arm64 amd64"

# build the app for each plat/arch
echo "Built executable: $exec"
for platform in $platforms; do
  for arch in $architectures; do
    echo "Building for $platform-$arch..."
    output="${exec}-${platform}-${arch}"
    GOOS=$platform GOARCH=$arch go build -o "$output" "$app"
  done
done
