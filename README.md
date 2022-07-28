# Rooms API

## Description

Simple multi-room chat application developed with [Gin](https://gin-gonic.com/), [GORM](https://gorm.io/index.html) and [Gorilla WebSocket](https://pkg.go.dev/github.com/gorilla/websocket).\
Work in progress...

## Environment Variables

Create a `.env` file in the root directory. And add these default values:

```
DATABASE_URL=postgres://postgres:password@postgres:5432/rooms
API_SECRET=SecretSecretSecret
TOKEN_HOUR_LIFESPAN=12
ENVIRONMENT=dev
```

## Run Application with Docker

More information about [Docker](https://www.docker.com/).\
To run the application type this command in the root folder.

```bash
$ docker compose up
```
