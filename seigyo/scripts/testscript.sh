#!/bin/sh

for i in `seq 1 10`; do
  curl http://frontend.test:8081;
done
