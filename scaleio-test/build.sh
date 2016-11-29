#! /usr/bin/env bash

rm -rf ./vendor
rm glide.lock
rm ./scaleio-test
glide up

grep -R --exclude-dir vendor --exclude-dir .git --exclude-dir mesos --exclude build.sh TODO ./

if [ "$1" == "native" ]; then
echo "Building Native Binary"
go build .
else
echo "Building Linux amd64 Binary"
GOOS=linux GOARCH=amd64 go build .
fi