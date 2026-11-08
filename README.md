# Neoroute

![Coverage](https://img.shields.io/badge/Coverage-58.4%25-yellow)

> [!WARNING]
> You are currently viewing the `v2.0.0` release candidate version of this repository. We may still ship some minor changes, but according to our testing this version should pretty much be stable. However, we're releasing it as release candidate since we've not verified this release in production yet. Major changes are listed at [https://liphium.dev/neoroute-v2](https://liphium.dev/neoroute-v2).

Neoroute is a **batteries-included remote procedure call (RPC) framework** for Golang, running exclusively on top of Web primitives (currently HTTP and WebSocket). With our rich tooling, you get **code generation**, an **interactive debugger** for sending requests and a lot more with minimal setup. We have everything you love about RPC and web frameworks, all in one **fully integrated package and ecosystem**.

## Installation

To use Neoroute you need Go version **1.27** or higher. Just add it to your project using `go`:

```bash
go get -u github.com/Liphium/neoroute/v2@latest
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

## Migrating to v2

To migrate to v2, there are basically no major steps except for changing all of your install paths to the proper v2 versions of the packages. Here is a list of the new install paths:

| Module/package | New install/import path |
| --- | --- |
| Core framework | `github.com/Liphium/neoroute/v2` |
| Neogen command | `github.com/Liphium/neoroute/v2/cmd/neogen` |
| Server HTTP transporter | `github.com/Liphium/neoroute/transporter/http/v2` |
| Server WebSocket transporter | `github.com/Liphium/neoroute/transporter/websocket/v2` |
| Go client SDK | `github.com/Liphium/neoroute/client/sdk/v2` |
| Client HTTP transporter | `github.com/Liphium/neoroute/client/transporter/http/v2` |
| Client WebSocket transporter | `github.com/Liphium/neoroute/client/transporter/websocket/v2` |
| Neodebug | `github.com/Liphium/neoroute/pkg/neodebug/v2` |
| Neodebug config | `github.com/Liphium/neoroute/pkg/neodebug/v2/config` |

Each Go module has its own `/v2` suffix. Packages within a module come after that suffix, as shown by Neogen and Neodebug config above.

You can find a list of the recent major changes at [https://liphium.dev/neoroute-v2](https://liphium.dev/neoroute-v2).

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
