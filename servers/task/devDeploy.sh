#!/bin/bash
docker build -t kangwooc/task .

docker push kangwooc/task
docker network disconnect finalnetwork tasking
docker rm -f tasking

docker run -d \
--name tasking \
--network finalnetwork \
-e MONGOADDR=mongo:27017 \
-e RABBITADDR=rabbitsvr:5672 \
kangwooc/task
