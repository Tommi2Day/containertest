#!/bin/bash

export DBUSER='c##test'
export DBPASS='Test!Password'
export TNS_ADMIN=$(pwd)
export ORACLE_PASSWORD='XE-Manager21'

docker run -d -p 1521:1521 \
  --env ORACLE_PASSWORD=$ORACLE_PASSWORD \
  --name oracle-xe -d gvenzl/oracle-xe:21.3.0-slim

echo  "Wait 60s for the database to be up and running"
sleep 60
sqlplus -l sys/$ORACLE_PASSWORD@XE as sysdba<<EOF
@create_user.sql
  exit;
EOF

go mod vendor
go run main.go xe xepdb1

docker stop oracle-xe
docker rm oracle-xe