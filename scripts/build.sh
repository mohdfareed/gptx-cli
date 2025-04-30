#!/bin/sh

app=./app
exec=.bin/chat

platforms="linux darwin windows"
architectures="arm64 amd64"

for platform in $platforms; do
  for arch in $architectures; do
    echo "Building for $platform-$arch..."
    output="${exec}-${platform}-${arch}"
    GOOS=$platform GOARCH=$arch go build -o "$output" "$app"
  done
done
