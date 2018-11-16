# Tags Drive Core

This repository contains backend part of **Tags Drive**

##

- [Tags Drive Core](#tags-drive-core)
  - [Security](#security)
  - [File structure](#file-structure)
    - [Config folder](#config-folder)
      - [JSON storage](#json-storage)
    - [Data folder](#data-folder)
    - [SSL folder](#ssl-folder)
  - [API](#api)
    - [Files](#files)
    - [Tags](#tags)
  - [Additional info](#additional-info)
    - [Environment variables](#environment-variables)

## Security

Uploaded files can be encrypted. **Tags Drive** uses sha256 sum of the password for encryption. Encryption is realized by [minio/sio](https://github.com/minio/sio) package.

## File structure

### Config folder

- `tokens.json` - contains valid tokens

  <details>
    <summary>Example</summary>

    ```json
    [
      {
        "token": "first-token",
        "expire": "2018-12-13T17:13:02.7716523+03:00"
      },
      {
        "token": "second-token",
        "expire": "2019-01-02T15:35:18.7829909-08:00"
      }
    ]
    ```
  </details>

#### JSON storage

- `files.json` - contains json map of all files

  <details>
    <summary>Example</summary>

    ```json
    {
      "1.jpg": {
        "filename": "1.jpg",
        "type": "image",
        "origin": "data/1.jpg",
        "description": "some cool image",
        "size": 527928,
        "tags": [12, 15, 17, 19, 18],
        "addTime": "2018-10-12T20:37:54.5515067+03:00",
        "preview": "data/resized/1.jpg"
      },
      "file.txt": {
        "filename": "file.txt",
        "type": "file",
        "origin": "data/file.txt",
        "description": "",
        "size": 48,
        "tags": [],
        "addTime": "2018-11-04T23:54:54.9669548-08:00"
      }
    }
    ```
  </details>

- `tags.json` - contains json map of all tags

  <details>
    <summary>Example</summary>

    ```json
      {
        "12": {
          "id": 12,
          "name": "cute",
          "color": "#55dcd4"
        },
        "15": {
          "id": 15,
          "name": "nature",
          "color": "#c9f898"
        }
      }
    ```
  </details>

### Data folder

Folder `data` is used as a file storage

### SSL folder

Folder `ssl` contains TLS certificate files `cert.cert`, `key.key`

Use this command to generate self-signed TLS certificate:

`openssl req -x509 -nodes -newkey rsa:2048 -sha256 -keyout key.key -out cert.cert`

## API

### Files

- `GET /api/files`

  **Params:**
  - **expr**: logical expression. Example: `!(12&15)&(12|15)` means all files with single tag with id `12` or `15`
  - **search**: text for search
  - **sort**: name | size | time
  - **order**: asc | desc

  **Response:** json array of:

  ```go
  type FileInfo struct {
      Filename    string    `json:"filename"`
      Type        string    `json:"type"`
      Origin      string    `json:"origin"`
      Description string    `json:"description"`
      Size        int64     `json:"size"`
      Tags        []int     `json:"tags"`
      AddTime     time.Time `json:"addTime"`
      Preview string `json:"preview,omitempty"` // It's not empty, if Type == "iamge"
  }
  ```

- `GET /api/files/recent`

  **Params:**
  - **number**: number of returned files (5 is a default value)

  **Response:** same as `GET /api/files`

- `GET /api/files/download`

  **Params:**
  - **file**: file for downloading (to download multiple files at a time, use `file` several times: `file=123.jp  file=hello.png`)

  **Response:** zip archive

- `POST /api/files`
  
  **Body** must be `multipart/form-data`

  **Response:** json array of:

  ```go
  type multiplyResponse struct {
      Filename string `json:"filename"`
      IsError  bool   `json:"isError"`
      Error    string `json:"error"`
      Status   string `json:"status"` // Status isn't empty when IsError == false
  }
  ```

- `POST /api/files/recover`

  **Params**:
  - file: file for recovering (to recover multiple files at a time, use `file` several times:`file=123.jpg&file=hello.png`)

  **Response**: -

- `PUT /api/files/name`

  **Params:**
  - **file**: old filename
  - **new-name**: new filename

  **Response:** -

- `PUT /api/files/tags`

  **Params:**
  - **file**: filename
  - **tags**: updated list of tags, separated by comma (`tags=1,2,3`)

  **Response:** -

- `PUT /api/files/description`

  **Params:**
  - **file**: filename
  - **description**: updated description

  **Response:** -

- `DELETE /api/files`

  **Params:**
  - **file**: file for deleting (to delete multiplefiles at a time, use `file` several times:`file=123.jpg&file=hello.png`)
  - **force**: should file be deleted right now (if it isn't empty, file will be deleted right now)

  **Response:** json array of:

  ```go
  type multiplyResponse struct {
      Filename string `json:"filename"`
      IsError  bool   `json:"isError"`
      Error    string `json:"error"`
      Status   string `json:"status"` // Status isn't empty when IsError == false
  }
  ```

### Tags

- `GET /api/tags`

  **Params:** -

  **Response:** json map `tagID: Tag`, where

  ```go
  type Tag struct {
      ID    int    `json:"id"`
      Name  string `json:"name"`
      Color string `json:"color"`
  }
  ```

- `POST /api/tags`

  **Params:**
  - **name**: name of a new tag
  - **color**: color of a new tag (`#ffffff` by default)

  **Response:** -

- `PUT /api/tags`

  **Params:**
  - **id**: id of a tag
  - **name**: new name of a tag (can be empty)
  - **color**: new color of a tag (can be empty)

  **Response:** -

- `DELETE /api/tags`

  **Params:**
  - **id**: id of a tag (one tag at a time)

  **Response:** -

## Additional info

### Environment variables

| Variable   | Default | Description                                      |
| ---------- | ------- | ------------------------------------------------ |
| PORT       | 80      | Port for website                                 |
| TLS        | true    | Should **Tags Drive** use https                  |
| LOGIN      | user    | Login for login                                  |
| PSWRD      | qwerty  | Password for login                               |
| ENCRYPT    | false   | Should the **Tags Drive** encrypt uploaded files |
| DBG        | false   |                                                  |
| SKIP_LOGIN | false   | Let use **Tags Drive** without loginning         |
