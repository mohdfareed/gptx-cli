#!/bin/sh

# USAGE:
#    scripts/release.sh version
# ARGS:
#    version - the release version tag, e.g. v0.0.0

usage="usage: $0 v0.0.0"
if [ "$#" -ne 1 ]; then echo "$usage" && exit 1; fi

echo "creating release: $1"
echo "press enter to continue"
read _

git tag $1
git push origin $1
