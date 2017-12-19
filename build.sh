#!/usr/bin/env bash

GOPATH=/Users/ian/Code/GO/workspace
BUILDPATH=$GOPATH/bin/tracr
MAINPATH=$GOPATH/src/tracr-daemon/main/main.go

OS=linux
PLATFORM=amd64

echo "Building executable for target architecture - "$OS"/"$PLATFORM
env GOOS=$OS GOARCH=$PLATFORM go build -i -o $BUILDPATH/tracrd $MAINPATH

echo ""
echo "Built executable"
echo ""
echo $BUILDPATH
ls -l $BUILDPATH

echo ""
echo "DONE"