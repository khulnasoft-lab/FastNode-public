#!/bin/bash

DEST="/deploy"
LOGDIR="/var/fastnode/log"

TARGET="mockserver"
HOST="mock.khulnasoft.com"

scp mockserver.py mock.khulnasoft.com:/deploy
ssh mock.khulnasoft.com "sudo killall gunicorn"
ssh mock.khulnasoft.com "sudo -b gunicorn --pythonpath /deploy --bind 0.0.0.0:80 mockserver:app"
