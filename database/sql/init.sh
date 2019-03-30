#!/bin/bash

psql --command "CREATE USER postgres WITH SUPERUSER PASSWORD 'postgres';"
createdb -O postgres postgres
psql -d postgres -c "CREATE EXTENSION IF NOT EXISTS citext;"
psql postgres -f ./init.sql