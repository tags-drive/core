#!/bin/sh

export DEBUG=true
export WEB_PORT=:80
export WEB_TLS=false
export WEB_LOGIN=user
export WEB_PASSWORD=qwerty
export WEB_SKIP_LOGIN=true
export STORAGE_ENCRYPT=false

./tags-drive
