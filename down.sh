#/bin/bash

cd $(dirname $0)

cd seigyo
docker-compose down

cd ..

cd taisyou
docker-compose down
