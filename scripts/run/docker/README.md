# scripts/run/docker

This script runs **Tags Drive** in a Docker container

## Launch command

```bash
python scripts/run/docker/run.py
```

## Available flags

| Flag               | Description                            | Default value  |
| ------------------ | -------------------------------------- | -------------- |
| `--image-name`     | name of a Docker image                 | dev-tags-drive |
| `--image-tag`      | tag of a Docker image                  | latest         |
| `--container-name` | name of a Docker container             | dev-tags-drive |
| `--container-port` | port of a Docker container             | 80             |
| `--mount-folder`   | folder for mount to a Docker container | root folder    |
| `--build-only`     | don't run s Docker container           | False          |

**Note:** also you can use `--help` to display all available flags
