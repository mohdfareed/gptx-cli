#!/bin/sh

USAGE="usage: $0 version
The version is the tag name to be created, e.g. v1.0.0"
if [ "$#" -ne 1 ]; then echo "$USAGE" && exit 1; fi

echo "creating release: $1"
echo "press enter to continue"
read _

git tag $1
git push origin $1
