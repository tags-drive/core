# Tags Drive

**Tags Drive** is an open source standalone drive. The main feature of **Tags Drive** is that files have flat structure (there's no folders). Instead, every file has a tag (or tags).

## Why should I prefer Tags Drive to other cloud drives

For example, you want to save an image of a cat. You can save it into folder `cats` or into folder `cute`. Of course, you may keeps 2 equal files, but it would be better to use tags system. So, you just need to add tags `cats` and `cute` to the photo.

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
