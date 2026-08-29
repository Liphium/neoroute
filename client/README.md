# Neoroute Go Client

This is the official library to interact with a server written with [Neoroute](https://github.com/Liphium/neoroute) in Golang.

## Installation

Run the following command in your Go project:

```sh
go get -u github.com/Liphium/neoroute/client@latest
```

Depending on which transporter you would like to use, run one of the following commands as well:

```sh name="HTTP" key="http"
go get -u github.com/Liphium/neoroute/client/transporter/http@latest
```

```sh name="WebSocket" key="ws"
go get -u github.com/Liphium/neoroute/client/transporter/websocket@latest
```

## Features

- Supports [WebSocket](https://liphium.dev/neoroute/client/websocket) and [HTTP](https://liphium.dev/neoroute/client/http) transporters (everything available in Neoroute as of writing this)
- Code generation support via [neogen](https://liphium.dev/neoroute/utility/neogen)
- Stable, used in lots of projects, as well as [neodebug](https://liphium.dev/neoroute/utility/neodebug) and other internal Liphium projects

## Documentation

We have official documentation available at [liphium.dev/neoroute](https://liphium.dev/neoroute). Check the "Client" category for all of the guides related to this SDK.

## Examples

There are multiple client examples in the main [Neoroute repository](https://github.com/Liphium/neoroute/tree/main/examples). All of the client examples in the repository are written in Go, and therefore use this SDK.

Guides for how to run them are available in the README of all the examples.
