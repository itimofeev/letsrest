#!/usr/bin/env bash

export GOPATH=/Users/ilyatimofee/prog/go/letsrest/

PROJECT_PATH=/Users/ilyatimofee/prog/go/letsrest/src/github.com/itimofeev/letsrest
FRONTEND_PROJECT_PATH=/Users/ilyatimofee/prog/js/letsrest-ui

rm ${PROJECT_PATH}/tools/target/*
mkdir ${PROJECT_PATH}/tools/target


export GOOS=linux
export GOARCH=amd64
go build -v github.com/itimofeev/letsrest/main/letsrest



docker build --force-rm=true -t letsrest -f ${PROJECT_PATH}/tools/letsrest.Dockerfile .
docker save -o "$PROJECT_PATH/tools/target/letsrest.img" "letsrest"


rm letsrest

cp ${PROJECT_PATH}/tools/prod-files/* ${PROJECT_PATH}/tools/target/

echo 'building frontend'

#npm build ${FRONTEND_PROJECT_PATH}
cp -r ${FRONTEND_PROJECT_PATH}/build ${PROJECT_PATH}/tools/target/frontend
cd tools/target
tar -jcvf ${PROJECT_PATH}/tools/target/frontend.tar.bz2 frontend

rm -r frontend
