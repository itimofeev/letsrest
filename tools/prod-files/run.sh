#!/usr/bin/env bash

docker load -i letsrest.img

tar -jxvf frontend.tar.bz2
chown ilyaufo frontend

docker-compose -p letsrest -f prod.docker-compose.yml up -d --build
