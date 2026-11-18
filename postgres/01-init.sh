#!/bin/sh
set -eu
# Every application role owns only its database; PUBLIC cannot connect.
for scope in users orders notifications; do
 psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres -v scope="$scope" -v password="${scope}_local" <<'SQL'
CREATE ROLE :"scope" LOGIN PASSWORD :'password';
CREATE DATABASE :"scope" OWNER :"scope";
REVOKE ALL ON DATABASE :"scope" FROM PUBLIC;
SQL
 psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$scope" -v scope="$scope" <<'SQL'
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO :"scope";
SQL
 PGUSER="$scope" PGPASSWORD="${scope}_local" psql -v ON_ERROR_STOP=1 --dbname "$scope" -f "/schemas/${scope}.sql"
done
