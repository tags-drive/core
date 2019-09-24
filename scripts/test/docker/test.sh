#!/bin/bash

docker-compose -f ./scripts/test/docker/docker-compose.yml \
    up \
    --build \
    --abort-on-container-exit \
    --exit-code-from tags-drive
