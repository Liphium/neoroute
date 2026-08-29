# Chat server example

This example shows how to build a basic chat server with Neoroute.

You can connect using the following client transporters:

- **HTTP:** Send messages to the server.
- **WebSocket:** Send messages and receive message events from the server.

It's really simple and only persists the messages in memory.

## Running the example

Literally just what you're used to:

```sh
go run .
```

## Debugging the server

[Neodebug](https://liphium.dev/neoroute/utility/neodebug) is integrated into this project as well. It allows you to open a TUI interface, sort of like an API client, that lets you check functionality of the server.

To start it for HTTP use:

```sh
go run . --debug-http
```

For WebSocket use:

```sh
go run . --debug-ws
```

## Testing

This server also has some basic tests in `main_test.go`. Run them like this:

```sh
go test ./...
```
