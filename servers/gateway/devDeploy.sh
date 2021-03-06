./build.sh

export TLSCERT=/tls/fullchain.pem
export TLSKEY=/tls/privkey.pem
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)

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

sleep 40

docker run -d \
--hostname rabbit \
--name rabbitsvr \
--network finalnetwork \
rabbitmq:3-management

sleep 40

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
-v $(pwd)/tls:/tls:ro \
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