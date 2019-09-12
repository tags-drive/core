# scripts/test/docker

## Launch command

```bash
docker-compose -f ./scripts/test/docker/docker-compose.yml \
    up \
    --build \
    --abort-on-container-exit \
    --exit-code-from tags-drive
```

**Params:**

- `-f scripts/test/docker/docker-compose.yml` – use a specail `docker-compose.yml` file
- `--build` – build an image every time
- `--abort-on-container-exit` – stop containers after finish of the tests
- `--exit-code-from tags-drive` – exit with the same code as `tags-drive` service
