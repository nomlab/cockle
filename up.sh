#/bin/bash

cd $(dirname $0)

cd taisyou
docker-compose up -d

cd ..
cd seigyo
docker-compose up -d
