# Neoroute

![Coverage](https://img.shields.io/badge/Coverage-49.6%25-yellow)

Neoroute is a **batteries-included remote procedure call (RPC) framework** for Golang, running exclusively on top of Web primitives (currently HTTP and WebSocket). With our rich tooling, you get **code generation**, an **interactive debugger** for sending requests and a lot more with minimal setup. We have everything you love about RPC and web frameworks, all in one **fully integrated package and ecosystem**.

## Installation

To use Neoroute you need Go version **1.27** or higher. Just add it to your project using `go`:

```bash
go get -u github.com/Liphium/neoroute@latest
```

## Features

- Use the same handlers over multiple protocols
- Rich [routing system](https://liphium.dev/neoroute/guides/routing) with support for middlewares, groups and more
- [Events & adapters](https://liphium.dev/neoroute/guides/events-adapters) for easy server to client communication
- Full support for [HTTP](https://liphium.dev/neoroute/guides/http) and [WebSocket](https://liphium.dev/neoroute/guides/websocket) transporters
- Fast and small messages thanks to [MessagePack](https://msgpack.org/)
  - Some of the fastest encoding and decoding thanks to [msgp](https://github.com/tinylib/msgp)'s code generation
- Client SDKs for [Go](https://liphium.dev/neoroute/sdk/go) and [TypeScript](/neoroute/sdk/typescript)
- Full code generation support for Go and TypeScript using [neogen](https://liphium.dev/neoroute/utility/neogen)
- Interactive debugging CLI using [neodebug](https://liphium.dev/neoroute/utility/neodebug)
- Lots of [testing utilities](https://liphium.dev/neoroute/guides/testing)

## Documentation

We have official documentation available at [liphium.dev/neoroute](https://liphium.dev/neoroute). Feel free to check it out.

There are also examples in the `examples` folder right here in this repository.

## Use with AI

We provide the entire documentation you're seeing right now for AI as well. It's index is available at [ai.liphium.dev/neoroute/index.md](https://ai.liphium.dev/neoroute/index.md).

You can install our `liphium-neoroute` skill that contains this information and some rules for Neoroute using:

```sh
npx skills add liphium/dev@liphium-neoroute
```

Your agent will love it.

## Contributing

We don't have any hard rules, but please keep in mind the following:

- No full AI PRs or issues: You should know what your code and writing contains, don't commit straight up slop.
- Be friendly with everyone, we're all here working on this project in our free time.

If you then want to contribute, follow this simple process:

1. Create an issue outlining what you want to add or the bug you found.
2. If it's a feature, discuss it with us and we'll see if we add anything like it. Specify if you want to work on it or have already started working on it.
3. If you chose contribution: Create a PR and iterate on it with us, we'll check and review when we have time.
4. When the feature / bug fix is out, the issue will be closed.

We really appreciate any contribution, even if it's just fixing a little thing or improving some documentation.

Thanks to everyone who contributed to Neoroute!

Read more at [liphium.dev/neoroute](https://liphium.dev/neoroute)...
