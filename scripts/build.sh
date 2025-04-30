#!/bin/sh

# USAGE:
#    scripts/build.sh [path]
# ARGS:
#    path - the build output path (default: ./.bin/chatgpt)

usage="usage: $0 [path=.bin/chatgpt]"
if [ "$#" -gt 1 ]; then echo "$usage" && exit 1; fi

app=./chatgpt # the app source code
exec="${1:-.bin/chatgpt}" # the built executable path

# the supported platforms and architectures
platforms="darwin linux windows"
architectures="arm64 amd64"

# build the app for each plat/arch
echo "built executable: $exec"
for platform in $platforms; do
  for arch in $architectures; do
    echo "building for $platform-$arch..."
    output="${exec}-${platform}-${arch}"
    GOOS=$platform GOARCH=$arch go build -o "$output" "$app"
  done
done
