# Tags Drive

**Tags Drive** is an open source standalone drive. The main feature of **Tags Drive** is that files have flat structure (there's no folders). Instead, every file has a tag (or tags).

## Why should I prefer Tags Drive to other cloud drives

For example, you want to save an image of a cat. You can save it into folder `cats` or into folder `cute`. Of course, you may keeps 2 equal files, but it would be better to use tags system. So, you just need to add tags `cats` and `cute` to the photo.

## API

All API methods require auth.

### Files

- `GET /api/files?sort=(name|size|time)&order(asc|desc)&tags=first,second,third&mode=(or|and|not)&search=abc` - get list of files.
- `GET /api/files/recent?number=5` - get list of the last uploaded files ((5 is a default number of returned files)
- `PUT /api/files?oldname=123&newname=567` - rename file.
- `POST /api/files` - upload files (`Content-Type: multipart/form-data`)
- `DELETE /api/files?file=file1&file=file2` - delete file.

```go
// multiplyResponse is used as response by POST /api/files and DELETE /api/files
type multiplyResponse struct {
	Filename string `json:"filename"`
	IsError  bool   `json:"isError"`
	Error    string `json:"error"`
	Status   string `json:"status"` // Status isn't empty when IsError == false
}
```

### Tags

- `PUT /api/tags?file=filename&tags=first,second,third` - change tags.

## Install

### Backend

**Requirements:**

- Docker
- Docker Compose

**Parameters:**

| Environment | Default | Description                     |
| ----------- | ------- | ------------------------------- |
| PORT        | 80      | Port for website                |
| TLS         | false   | Should **Tags Drive** use https |
| LOGIN       | user    | Login for login                 |
| PSWRD       | qwerty  | Password for login              |
| DBG         | false   |                                 |

### Frontend

**Tags Drive** uses framework [Vue.js](https://vuejs.org).
