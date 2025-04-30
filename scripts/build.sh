#!/bin/sh

app=./chat # the app source code
exec=.bin/chat # the built executable path

# build for the current platform
echo "Building app..."
go build -o "$exec" "$app"

# the supported platforms and architectures
platforms="darwin linux windows"
architectures="arm64 amd64"

# build the app for each plat/arch
for platform in $platforms; do
  for arch in $architectures; do
    echo "Building for $platform-$arch..."
    output="${exec}-${platform}-${arch}"
    GOOS=$platform GOARCH=$arch go build -o "$output" "$app"
  done
done
