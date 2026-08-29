# Chat client HTTP example

This example needs the [chat server example](https://github.com/Liphium/neoroute/tree/main/examples/chat_server) to be running in order to work. It shows you how you can connect to the server using the HTTP transporter from the Go client SDK.

## Running the example

Just simply use this command:

```sh
go run .
```

## Generating definitions

A command you may not be familiar with:

```sh
go generate ./...
```

This will use [neogen](https://liphium.dev/neoroute/utility/neogen) to generate the definitions from the server. You can use them to more easily connect.

You can check the `generated` folder to see what this generated code looks like.
