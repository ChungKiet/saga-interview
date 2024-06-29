# WebSocket Server with Unique Big Integer Responses

This repository contains a Go application that sets up a WebSocket server. The server accepts WebSocket connections and responds to any message with a unique random big integer. The uniqueness of the integers is ensured by storing and checking the numbers in a Redis database.

## Requirements

- Go 1.20+
- Docker (for running Redis)
- A WebSocket client (for testing)

## Setup

### 1. Run Redis in Docker

First, ensure you have Docker installed on your machine. Then, run the following command to start a Redis server in a Docker container:

```sh
docker run --name redis-server -p 6379:6379 -d redis
```

### 2. Install Go Dependencies

Run the following commands to get the required Go packages:

```sh
go get github.com/go-redis/redis/v8
go get github.com/gorilla/websocket
```

### 3. Run the Go Application
```sh
go run main.go
```

### 4. Test the WebSocket Server

You can use a WebSocket client to connect to `ws://localhost:8080/ws` and send any message. The server should respond with a unique random big integer.

You also can use file index.html to test web socket with multi connections.

