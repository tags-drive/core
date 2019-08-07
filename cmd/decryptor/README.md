# Decryptor

Decryptor decrypts files encrypted by **Tags Drive**

## Usage

1. CD to **Tags Drive** root folder (there must be the `var` folder)
2. Create a folder `decrypted-files`
3. Run

    ```bash
    docker run --rm \
        -v $PWD/var:/app/var \
        -v $PWD/decrypted-files:/app/ecrypted-files \
        kirtis/tags-drive decrypt --phrase=PHRASE
    ```

4. Decrypted files will be saved in the `decrypted-files` folder

### CL args

| Arg                     | Default             | Description                       | Required |
| ----------------------- | ------------------- | --------------------------------- | -------- |
| `--phrase`              |                     | Pass phrase used to encrypt files | yes      |
| `--config-file`         | `./var/files.json`  | Path to JSON file with files info |          |
| `--data-folder`         | `./var/data`        | Path to files                     |          |
| `-o`, `--output-folder` | `./decrypted-files` | Folder to save decrypted files    |          |
