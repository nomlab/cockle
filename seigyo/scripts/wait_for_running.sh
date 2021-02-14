#!/bin/sh


while ! curl http://frontend.test:8081/health 
do
    echo "waiting"
    sleep 1
done


while ! curl http://backend.test:9000/health 
do
    echo "waiting"
    sleep 1
done
