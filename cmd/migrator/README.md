# Migrator

Migrator migrates files from **Disk** to **S3 Storage** and vice versa.

## Usage

1. CD to **Tags Drive** root folder (there must be the `var` folder)
2. Run **Migrator**. Example:

    ```bash
    docker run --rm \
        -v $PWD/var:/app/var \
        kirtis/tags-drive migrate \
        --from=disk \
        --to=s3 \
        --disk.encrypted \
        --disk.pass-phrase=some_pass_phrase \
        --s3.endpoint=127.0.0.1:9000 \
        --s3.access-key=login \
        --s3.secret-key=password \
        --s3.secure
    ```

### CL args

| Arg                    | Default | Options      | Required |
| ---------------------- | ------- | ------------ | -------- |
| `--from`               |         | `disk`, `s3` | yes      |
| `--to`                 |         | `disk`, `s3` | yes      |
| `--disk.encrypted`     | `false` |              |          |
| `--disk.pass-phrase`   |         |              |          |
| `--s3.endpoint`        |         |              | yes      |
| `--s3.access-key`      |         |              | yes      |
| `--s3.secret-key`      |         |              | yes      |
| `--s3.secure`          | `false` |              |          |
| `--s3.bucket-location` |         |              |          |
