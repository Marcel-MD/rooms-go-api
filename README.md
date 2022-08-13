# Rooms API

## Description

Multi-room chat application developed with [Gin](https://gin-gonic.com/), [GORM](https://gorm.io/index.html), [Go Redis](https://redis.uptrace.dev/) and [Gorilla WebSocket](https://pkg.go.dev/github.com/gorilla/websocket).
Work in progress...

## Environment Variables

Create a `.env` file in the root directory. And add these default values:

```
DATABASE_URL=postgres://postgres:password@postgres:5432/rooms
REDIS_URL=redis://:password@redis:6379/0
API_SECRET=SecretSecretSecret
TOKEN_HOUR_LIFESPAN=12
ENVIRONMENT=dev
```

## Run Application with Docker

More information about [Docker](https://www.docker.com/).
To run the application type this command in the root folder.

```bash
$ docker compose up
```

You might have to run this command twice if it doesn't work the first time :)

## API Endpoints

For authentication are used bearer tokens.

- **User** `/api/users`

  - [GET] `/` - Get all users

  - [GET] `/:id` - Get user by ID

  - [POST] `/register` - Register User

    ```json
    {
      "firstName": "First",
      "lastName": "Last",
      "email": "firstlast@mail.com",
      "password": "password"
    }
    ```

  - [POST] `/login` - Login User

    ```json
    {
      "email": "firstlast@mail.com",
      "password": "password"
    }
    ```

- **Room** `/api/rooms`

  - [GET] `/` - Get all rooms

  - [GET] `/:id` - Get room by ID

  - [POST] `/` - Create room

    ```json
    {
      "name": "room"
    }
    ```

  - [PUT] `/:id` - Update room by ID

    ```json
    {
      "name": "updated room"
    }
    ```

  - [DELETE] `/:id` - Delete room by ID

  - [POST] `/:room_id/users/:user_id` - Add user to room

  - [DELETE] `/:room_id/users/:user_id` - Remove user from room

- **Message** `/api/messages`

  - [GET] `/:room_id?page=1&size=10` - Get paginated messages by room ID

  - [POST] `/:room_id` - Create message

    ```json
    {
      "text": "Hello World!"
    }
    ```

  - [PUT] `/:id` - Update message by ID

    ```json
    {
      "text": "Goodbye World!"
    }
    ```

  - [DELETE] `/:id` - Delete message by ID

- **WebSocket** `/api/ws`

  - [GET] `/` - Connect to all user's rooms

## WebSocket

- Create Message

  ```json
  {
    "text": "Hello World!",
    "command": "CreateMessage",
    "targetId": "room_id",
    "roomId": "room_id"
  }
  ```

- Update Message

  ```json
  {
    "text": "Goodbye World!",
    "command": "UpdateMessage",
    "targetId": "message_id",
    "roomId": "room_id"
  }
  ```

- Delete Message

  ```json
  {
    "text": "anything",
    "command": "DeleteMessage",
    "targetId": "message_id",
    "roomId": "room_id"
  }
  ```

- Add User to Room

  ```json
  {
    "text": "anything",
    "command": "AddUser",
    "targetId": "user_id",
    "roomId": "room_id"
  }
  ```

- Remove User from Room

  ```json
  {
    "text": "anything",
    "command": "RemoveUser",
    "targetId": "user_id",
    "roomId": "room_id"
  }
  ```

- Create Room

  ```json
  {
    "text": "Room Name",
    "command": "CreateRoom",
    "targetId": "anything",
    "roomId": "anything"
  }
  ```

- Update Room

  ```json
  {
    "text": "New Room Name",
    "command": "UpdateRoom",
    "targetId": "room_id",
    "roomId": "room_id"
  }
  ```

- Delete Room

  ```json
  {
    "text": "anything",
    "command": "DeleteRoom",
    "targetId": "room_id",
    "roomId": "room_id"
  }
  ```
