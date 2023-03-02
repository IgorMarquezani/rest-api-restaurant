#!/bin/bash

if [[ $1 == "start" ]]; then
  docker-compose up -d
  go run ../main.go
elif [[ $1 == "reset" ]]; then
  docker-compose down
  for ((i = 1; i < 10; i++)); do
    container=$(docker image ls | sed '2,800 ! d' | sed -n {"$i p"} | awk {'print $1'} )
    if [[ $container == "postgres" || $container == "dpage/pgadmin4" ]]; then
      docker image rmi -f $container
      $i=$(($i-1))
    fi
  done
  docker volume rm $(docker volume ls -q) 2> /dev/null 
elif [[ $1 == "help" ]]; then
  if [[ $LANG == "pt_BR.UTF-8" ]]; then
    echo "help:  show this message"
    echo "start: start the containers"
    echo "reset: completely reset containers, images and volumes"
  fi
  else
  echo "help:  mostra esta mensagem"
  echo "start: inicia os containers"
  echo "reset: reinicia totalmente containers, imagens e volumes"
fi
