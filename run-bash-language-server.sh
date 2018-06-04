#!/bin/bash

# This script runs 3 services and exits if any of them die:
#
# - mongo (for explainshell)
# - explainshell web server
# - lsp-adapter for the bash-language server
#
# This script assumes that:
#
# - The mongo db is at /data/db2 (instead of /data/db, because /data/db in the
#   mongo image is read-only)

trap "exit" INT TERM
trap "kill 0" EXIT

# Spawn services in the background
mongod --dbpath /data/db2 &
# Listen on 0.0.0.0 so that it can accept connections from outside the Docker
# container
(cd explainshell && env HOST_IP=0.0.0.0 make serve) &
env EXPLAINSHELL_ENDPOINT=http://localhost:5000 lsp-adapter --trace --glob="*.sh:*.bash:*.zsh" --proxyAddress=0.0.0.0:8080 bash-language-server start &

# Exit if any of the services die
while sleep 1; do
  if ! pgrep mongod >/dev/null; then
    echo "ERROR mongod died"
    exit 1
  fi
  if ! pgrep make >/dev/null; then
    echo "ERROR make died"
    exit 1
  fi
  if ! pgrep lsp-adapter >/dev/null; then
    echo "ERROR lsp-adapter died"
    exit 1
  fi
done
