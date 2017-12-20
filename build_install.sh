#!/usr/bin/env bash

./build.sh
retVal=$?
if [ ! $retVal -eq 0 ]; then
    echo "Build failed - skipping installation"
    exit $retVal
fi
./install.sh
