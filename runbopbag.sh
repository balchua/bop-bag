#!/bin/bash

name=$NODE_NAME
service=$HEADLESS_SVC
nodeName=`hostname`
echo "Starting the node $nodeName"
sleep 5
if [[ $nodeName == *-0 ]] # * is used for pattern matching
then
  echo "Starting seed node. $nodeName"
  /app/bopbag  serve --db /data --certs /app/certs  --dbAddress "$nodeName:9000"
else
  main=`echo ${nodeName%-*}-0`
  echo "Joining node $nodeName to $main:9000"

  /app/bopbag serve --db /data --certs /app/certs --dbAddress "$nodeName:9000" --join "$main:9000"
fi
