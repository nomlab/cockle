#!/bin/bash

PID=$1
ADDR=$2

WORKING_DIR=$(dirname $0)
cd $WORKING_DIR

rm -rf image
./phaul-client $PID $ADDR
