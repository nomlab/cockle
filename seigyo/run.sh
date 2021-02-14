#!/bin/bash -e

ADDRc=$(docker inspect --format '{{ .NetworkSettings.Networks.phaulincontainer_phaul_test.IPAddress }}' frontend)
ADDRs=$(docker inspect --format '{{ .NetworkSettings.Networks.phaulincontainer_phaul_test.IPAddress }}' backend)

echo "s"
echo $ADDRs
echo "r"
echo $ADDRc

echo 'python /root/main.py --sip ' $ADDRs ' --cip ' $ADDRc ' --servicename frontend'
#docker run -it --rm --net="phaulincontainer_phaul_test" teiansys python3 /root/main.py --cpi=$ADDRc --spi=$ADDRs
docker run -it --rm -v /home/vagrant/teiansys/scripts:/root/scripts --net="phaulincontainer_phaul_test" --dns="192.168.16.3" teiansys /bin/bash


