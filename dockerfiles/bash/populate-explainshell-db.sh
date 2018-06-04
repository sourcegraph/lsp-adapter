#!/bin/bash

# This script is intended to be used during the build of the bash language
# server to initialize the explainshell database with a few thousand man pages.
# This is a pretty slow process which takes ~30 minutes.
#
# This script assumes:
#
# - The current working directory is explainshell
# - Everything that explainshell needs is installed (mongo, python, etc.)
# - The mongo db is at /data/db2 (instead of /data/db, because that's read-only
#   in the mongo image)

wait_for_mongo_up() {
  while true; do
    nc -zvv localhost 27017 && return
    sleep 1
  done
}

wait_for_server_up() {
  while true; do
    nc -zvv localhost 5000 && return
    sleep 1
  done
}

# Start mongo and the explainshell web server in the background
mongod --dbpath /data/db2 &
make serve &
wait_for_mongo_up
wait_for_server_up

# Load the classifiers for flags in man pages
mongorestore -d explainshell dump/explainshell && mongorestore -d explainshell_tests dump/explainshell

# Sanity check
make tests

# Download a few thousand man pages
git clone https://github.com/idank/explainshell-manpages
cd explainshell-manpages
git checkout bb7f4dfb037b890de58e0541b369cad1eb6ae07f
# Avoid busting the cache if newer commits get pushed
rm -rf .git
cd ..

# Actually load the man pages. Last I checked there were ~26,000 man pages and
# it took ~30 minutes to load. It loads in batches for speed, and the batch size
# is capped at 1000 to avoid hitting the limit on number of shell arguments.
find . -name "*.gz" | xargs -n 1000 env PYTHONPATH=. python explainshell/manager.py

# Done loading, can stop mongo and the explainshell web server.
jobs -p | xargs kill
