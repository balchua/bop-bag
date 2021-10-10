#!/bin/bash

name=$NODE_NAME
service=$HEADLESS_SVC
echo "Starting the node $name"
sleep 5
if [[ $name == *-0 ]] # * is used for pattern matching
then
  echo "Starting seed node. $name"
  /app/bopbag  serve --db /data --certs /app/certs  --dbAddress "$name.$service:9000"
else
  main=`echo ${name%-*}-0`
  echo "Joining node $name.$service to $main.$service:9000"

  /app/bopbag serve --db /data --certs /app/certs --dbAddress "$name.$service:9000" --join "$main.$service:9000"
fi
