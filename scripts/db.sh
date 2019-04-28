#!/bin/bash

# fail on error
set -e

# =============================================================================================
if [[ "$(basename $PWD)" == "scripts" ]]; then
	cd ..
fi
echo $PWD

# =============================================================================================
source .env
source ~/.config/ir.conf || true

# =============================================================================================
echo "waiting on postgres ..."
until PGPASSWORD=dev-secret psql -h 127.0.0.1 -U dev-user -d ircollector_db -c '\q'; do
  echo "waiting ..."
  sleep 2
done
echo "postgres is up!"

# =============================================================================================
echo "setting up ircollector_db ..."

# TODO: db migrations needed?
