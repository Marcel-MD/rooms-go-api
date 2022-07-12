# Rooms API

## Description

Simple multi-room chat application developed with [Gin](https://gin-gonic.com/), [GORM](https://gorm.io/index.html) and [Gorilla WebSocket](https://pkg.go.dev/github.com/gorilla/websocket).\
Work in progress...

## Environment Variable

Create a `.env` file in the root directory. And add these default values:

```
DATABASE_URL=postgres://postgres:password@localhost:5432/rooms
API_SECRET=SecretSecretSecret
TOKEN_HOUR_LIFESPAN=12
ENVIRONMENT=dev
```

## Docker PostgreSQL Setup

More information about [Docker](https://www.docker.com/).

```bash
$ docker-compose up
```

Use [pgAdmin](https://www.pgadmin.org/) to connect to PostgreSQL and create a new database named `rooms`.

## Installation

```bash
$ go mod tidy
```

## Running the app

```bash
$ go run .
```
