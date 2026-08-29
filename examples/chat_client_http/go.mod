module github.com/Liphium/neoroute/examples/chat_client_http

go 1.27.0

replace github.com/Liphium/neoroute/client => ../../client

replace github.com/Liphium/neoroute/client/transporter/http => ../../client/transporter/http

require (
	github.com/Liphium/neoroute/client v0.0.0
	github.com/Liphium/neoroute/client/transporter/http v0.0.0
	github.com/tinylib/msgp v1.6.4
)

require (
	github.com/philhofer/fwd v1.2.0 // indirect
)
