# Tags Drive Core

[![Go Report Card](https://goreportcard.com/badge/github.com/tags-drive/core)](https://goreportcard.com/report/github.com/tags-drive/core)

This repository contains the backend part of **Tags Drive**

##

- [Usage](#usage)
  - [Environment variables](#environment-variables)
- [Development](#development)
- [File structure](#file-structure)
  - [Var folder](#var-folder)
  - [SSL folder](#ssl-folder)
- [API](#api)
  - [General endpoints](#general-endpoints)
  - [Auth](#auth)
  - [Files](#files)
  - [Tags](#tags)
  - [Share](#share)
- [Additional info](#additional-info)
  - [Security](#security)
  - [General structures](#general-structures)

## Usage

### Environment variables

| Variable             | Default | Description                                                          |
| -------------------- | ------- | -------------------------------------------------------------------- |
| PORT                 | 80      | Port for http server                                                 |
| TLS                  | true    | Enable HTTPS                                                         |
| LOGIN                | user    | Set your login                                                       |
| PSWRD                | qwerty  | Set your password                                                    |
| SKIP_LOGIN           | false   | Skip the log-in procedure                                            |
| MAX_TOKEN_LIFE       | 1440h   | The max lifetime of a token (default lifetime is 60 days)            |
| ENCRYPT              | false   | Encrypt uploaded files                                               |
| PASS_PHRASE          | ""      | A phrase for file encryption. Cannot be empty if `ENCRYPT == true`   |
| TIME_BEFORE_DELETING | 168h    | Time before deleting a file from the Trash (default delay is 7 days) |
| DBG                  | false   |                                                                      |

## Development

There are two Python scripts that you can use to run a local version of the backend part:

- [scripts/run/run.py](scripts/run/run.py) – run a local version with `go run`. You can set env vars by editing the [.env file](scripts/run/run.env). It is the fastest way to launch the local version, but you need to have Go installed.
- [scripts/docker/run_docker.py](scripts/docker/run_docker.py) – build a Docker image and run a container. There are some command-line args (run `python scripts/docker/run_docker.py --help` to show all args)

## File structure

### Var folder

- `auth_tokens.json` - contains valid tokens

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

- `files.json` - contains a json map of all files

  <details>

    <summary>Example</summary>

    ```json
    {
      "1": {
        "id": 1,
        "filename": "cute-cat.jpg",
        "type": {
          "ext": ".jpg",
          "fileType": "image",
          "supported": true,
          "previewType": "image"
        },
        "origin": "data/1",
        "preview": "data/resized/1",
        "tags": [24,26],
        "description": "very cute cat :)",
        "size": 480900,
        "addTime": "2018-12-29T16:45:07.4440863+03:00",
        "deleted": false,
        "timeToDelete": "0001-01-01T00:00:00Z"
      },
    }
    ```

  </details>

- `tags.json` - contains a json map of all tags

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

#### Data folder

The `./var/data` folder is used as a file storage.

The `./var/data/resized` folder is used to store resized image previews.

### SSL folder

The `ssl` folder contains TLS certificate files `cert.cert` and `key.key`

Use this command to generate self-signed TLS certificate:

`openssl req -x509 -nodes -newkey rsa:2048 -sha256 -keyout key.key -out cert.cert`

## API

### General endpoints

- `GET /` – main page
- `GET /mobile` – mobile version
- `GET /share?shareToken=token` - **Tags Drive** in share mode
- `GET /login` – login page
- `GET /version` – returns version of the backend part
- `GET /data/{id}` – returns a file

### Auth

- `GET /api/user` – check if a user is authorized

  **Responses:**

  - if a user is authorized: `{ "authorized" : true }`
  - else: `401 Unauthorized`

- `POST /api/login` – sets cookie with auth token

  **Params:**
  - **login**: user's login
  - **password**: password (sha256 checksum repeated 11 times)

  **Response:** -

- `POST /api/logout` – deletes auth cookie

  **Params:** -

  **Response:** -

### Files

- `GET /api/file/{id}`

  **Params:**
  - **id**: id of a file
  - **shareToken** (optional): allow to use this API method without auth (the response (files, tags) can be limited)

  **Response:** json object of [`FileInfo`](#fileinfo)

- `GET /api/files`

  **Params:**
  - **expr**: logical expression. Example: `!(12&15)&(12|15)` means all files that have single tag with the id `12` or `15`
  - **search**: a text/regexp search
  - **regexp**: enable regexp search (it is `true` when **regexp** param is not an empty string)
  - **sort**: name | size | time
  - **order**: asc | desc
  - **offset**: lower bound `[offset:]`
  - **count**: number of returned files (`[offset:offset+count]`). If count == 0, all files will be returned. Default value is 0
  - **shareToken** (optional): allow to use this API method without auth (the response (files, tags) can be limited)

  **Response:** json array of [`FileInfo`](#fileinfo). Status code is `204` when offset is out of bounds.

- `GET /api/files/recent`

  **Params:**
  - **number**: number of returned files (5 is the default value)

  **Response:** json array of [`FileInfo`](#fileinfo)

- `GET /api/files/download`

  **Params:**
  - **ids**: list of file ids for downloading separated by commas `ids=1,2,54,9`
  - **shareToken** (optional): allow to use this API method without auth (the response (files, tags) can be limited)

  **Response:** zip archive

- `POST /api/files`
  
  **Params:**
  - **tags**: tags: list of tags separated by commas (`tags=1,2,3`)

  **Body** must be `multipart/form-data`

  **Response:** json array of [`multiplyResponse`](#multiplyresponse)

#### Changing file info

- `PUT /api/file/{id}/name`

  **Params:**
  - **id**: file id
  - **new-name**: new filename

  **Response:** updated file (json object of [`FileInfo`](#fileinfo))

- `PUT /api/file/{id}/tags`

  **Params:**
  - **id**: file id
  - **tags**: updated list of tags separated by commas (`tags=1,2,3`)

  **Response:** updated file (json object of [`FileInfo`](#fileinfo))

- `PUT /api/file/{id}/description`

  **Params:**
  - **id**: file id
  - **description**: updated description

  **Response:** updated file (json object of [`FileInfo`](#fileinfo))

#### Editing tags of multiple files

- `POST /api/files/tags`

  **Params:**
  - **files**: file ids (list of ids separated by ',')
  - **tags**: tags for adding (list of tags ids separated by ',')

  **Response:** -

- `DELETE /api/files/tags`

  **Params:**
  - **files**: file ids (list of ids separated by ',')
  - **tags**: tags for deleting (list of tags ids separated by ',')

  **Response:** -

#### Removing and recovering

- `DELETE /api/files`

  **Params:**
  - **ids**: list of file ids for deleting separated by commas `ids=1,2,54,9`
  - **force**: enable instant deletion (if not empty, files will be deleted immediately)

  **Response:** json array of [`multiplyResponse`](#multiplyresponse)

- `POST /api/files/recover`

  **Params**:
  - **ids**: list of file ids for recovering (list of ids separated by comma `ids=1,2,54,9`)

  **Response**: -

### Tags

- `GET /api/tags`

  **Params:**
  - **shareToken** (optional): allow to use this API method without auth (the response (files, tags) can be limited)

  **Response:** json object of [`Tags`](#Tag)

- `POST /api/tags`

  **Params:**
  - **name**: name of a new tag
  - **color**: color of a new tag (`#ffffff` by default)
  - **group**: group of a new tag (empty by default)

  **Response:** -

- `PUT /api/tag/{id}`

  **Params:**
  - **id**: tag id
  - **name**: new tag name (can be empty)
  - **color**: new tag color (can be empty)
  - **group**: new tag group (can be empty)

  **Response:** updated tag (json object of [`Tag`](#Tag))

- `DELETE /api/tags`

  **Params:**
  - **id**: tag id (one tag at a time)

  **Response:** -

### Share

- `GET /api/share/tokens` - returns all share tokens

  **Params:** -

  **Response:** json map with tokens and ids of shared files

  ```json
    {
      "token1": [1, 2, 3],
      "token2": [4, 5, 25],
    }
  ```

- `GET /api/share/token/{token}` - returns ids of files shared by passed token

  **Params:**
  - **token**: share token

  **Response:** json array with ids of shared files

- `POST /api/share/token` - create a new share token

  **Params:**
  - **ids**: list of ids of files to share separated by commas (example: `?ids=1,2,3`)

  **Response**: returns new share token

    ```json
      { "token": "created token" }
    ```

- `DELETE /api/share/token/{token}` - delete a share token

  **Params:**
  - **token**: share token

  **Response:** -

## Additional info

### Security

Uploaded files can be encrypted. **Tags Drive** uses sha256 sum of the `PASS_PHRASE` for encryption. Encryption is realized by [minio/sio](https://github.com/minio/sio) package.

### General structures

#### FileInfo

```go
    type FileType string

    type PreviewType string

    // Ext is a struct which contains a type of the original file and a type for a preview
    type Ext struct {
      Ext         string      `json:"ext"`
      FileType    FileType    `json:"fileType"`
      Supported   bool        `json:"supported"`
      PreviewType PreviewType `json:"previewType"`
    }

    type File struct {
      ID       int    `json:"id"`
      Filename string `json:"filename"`
      Type     Ext    `json:"type"`
      Origin   string `json:"origin"`
      Preview  string `json:"preview,omitempty"`
      //
      Tags        []int     `json:"tags"`
      Description string    `json:"description"`
      Size        int64     `json:"size"`
      AddTime     time.Time `json:"addTime"`
      //
      Deleted      bool      `json:"deleted"`
      TimeToDelete time.Time `json:"timeToDelete"`
    }
```

#### Tag

```go
  type Tag struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Color string `json:"color"`
    Group string `json:"group"`
  }

  type Tags map[int]Tag
```

#### multiplyResponse

```go
  type multiplyResponse struct {
      Filename string `json:"filename"`
      IsError  bool   `json:"isError"`
      Error    string `json:"error"`
      Status   string `json:"status"` // Status isn't empty when IsError == false
  }
```
