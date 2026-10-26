module github.com/Liphium/neoroute/examples/chat_client_ws

go 1.27.0

replace github.com/Liphium/neoroute/client/sdk/v2 => ../../client/sdk

replace github.com/Liphium/neoroute/client/transporter/websocket/v2 => ../../client/transporter/websocket

replace github.com/Liphium/neoroute/client/transporter/http/v2 => ../../client/transporter/http

require (
	github.com/Liphium/neoroute/client/sdk/v2 v2.0.0
	github.com/Liphium/neoroute/client/transporter/http/v2 v2.0.0
	github.com/Liphium/neoroute/client/transporter/websocket/v2 v2.0.0
	github.com/tinylib/msgp v1.6.4
)

require (
	github.com/coder/websocket v1.8.15 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
)
