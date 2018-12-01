#!/bin/bash
./build.sh

docker push kangwooc/final
docker push kangwooc/finaldb
# set the environment variable of "TLSCERT", "TLSKEY" and "MYSQL_ROOT_PASSWORD"
export TLSCERT=/etc/letsencrypt/live/api.kangwoo.tech/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.kangwoo.tech/privkey.pem
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)

ssh -i ~/.ssh/finalproject.pem ec2-user@52.33.171.173 "bash -s" << EOF
docker network disconnect finalnetwork redisserver
docker network disconnect finalnetwork mysqlserver
docker network disconnect finalnetwork mongo
docker network disconnect finalnetwork rabbitsvr
docker network disconnect finalnetwork gateway
docker network disconnect finalnetwork tasking
docker network rm finalnetwork

docker rm -f redisserver
docker rm -f mysqlserver
docker rm -f mongo
docker rm -f rabbitsvr
docker rm -f gateway
docker rm -f tasking

docker network create finalnetwork

docker pull kangwooc/finaldb
docker pull kangwooc/final

docker run -d \
--name redisserver \
--network finalnetwork \
redis

docker run -d --name mysqlserver \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=userDB \
--network finalnetwork \
kangwooc/finaldb

docker run -d \
--name mongo \
--network finalnetwork \
mongo

docker run -d \
--hostname rabbit \
--name rabbitsvr \
--network finalnetwork \
rabbitmq:3-management

sleep 20

docker run -d \
--name tasking \
--network finalnetwork \
-e MONGOADDR=mongo:27017 \
-e RABBITADDR=rabbitsvr:5672 \
kangwooc/task

docker run -d \
--name gateway \
--network finalnetwork \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e SESSIONKEY=$SESSIONKEY \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e REDISADDR=redisserver:6379 \
-e DBADDR=mysqlserver:3306 \
-e RABBITADDR=rabbitsvr:5672 \
-e SUMMARYADDR=summary:80 \
-e TASKADDR=tasking:80 \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
kangwooc/final

EOF