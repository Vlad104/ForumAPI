#!/bin/bash

psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';"
createdb -O docker docker
psql -d docker -c "CREATE EXTENSION IF NOT EXISTS citext;"
psql docker -f ./init.sql