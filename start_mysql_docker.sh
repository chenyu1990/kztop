#!/bin/bash

MYSQL_VERSION="8"
MYSQL_IMAGE=mysql:$MYSQL_VERSION

MYSQL_CONTAINER_NAME=mysql-$MYSQL_VERSION
HOST_MYSQL_CONF_ROOT=/srv/mysql/conf
HOST_MYSQL_DATA_ROOT=/srv/mysql/data

docker run -d \
  -p 33060:3306 \
  --name $MYSQL_CONTAINER_NAME \
  -v $HOST_MYSQL_DATA_ROOT:/var/lib/mysql \
  --restart always \
  -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
  $MYSQL_IMAGE

# copy mysql config files to host...
docker cp "$MYSQL_CONTAINER_NAME":/etc/mysql $HOST_MYSQL_CONF_ROOT
docker rm -f $MYSQL_CONTAINER_NAME
docker exec

docker run -d \
  -p 33060:3306 \
  --name $MYSQL_CONTAINER_NAME \
  -v $HOST_MYSQL_CONF_ROOT:/etc/mysql \
  -v $HOST_MYSQL_DATA_ROOT:/var/lib/mysql \
  --restart always \
  -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
  $MYSQL_IMAGE

docker exec mysql-8 sh -c "mysql -uroot -e \"CREATE USER IF NOT EXISTS 'kztop'@'172.%';GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP ON kztop.* TO 'kztop'@'172.%';flush privileges;CREATE DATABASE IF NOT EXISTS \`kztop\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;\""

#mysql -uroot -P 33060 -h 0.0.0.0 << EOF
#CREATE DATABASE IF NOT EXISTS \`kztop\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
#EOF