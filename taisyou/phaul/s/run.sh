#!/bin/bash

WORKING_DIR=$(dirname $0)
cd $WORKING_DIR

echo $(pwd)
rm -rf /tmp/livemig image
./phaul-server
