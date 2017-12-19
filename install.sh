#!/usr/bin/env bash

GOPATH=/Users/ian/Code/GO/workspace
BUILDPATH=$GOPATH/bin/tracr

echo "Installing..."
echo ""

echo "Transferring executable to home folder"
scp $BUILDPATH/tracrd iandpierce@35.196.123.75:~/
echo ""

echo "Installing executable"
ssh -t iandpierce@35.196.123.75 "sudo mv ~/tracrd /usr/local/bin"
echo ""

echo "DONE"
echo ""