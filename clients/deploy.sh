#!/bin/bash
./build.sh

docker push kangwooc/finalclient
ssh -i ~/.ssh/finalproject.pem ec2-user@35.165.176.110 "bash -s" << EOF
docker rm -f finalclient

docker pull kangwooc/finalclient
docker run -d \
--name finalclient \
-p 443:443 \
-p 80:80 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
kangwooc/finalclient

EOF