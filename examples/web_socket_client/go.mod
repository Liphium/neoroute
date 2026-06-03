module web_socket

go 1.26.3

replace github.com/Liphium/neoroute/client => ../../client/

require (
	github.com/Liphium/neoroute/client v0.0.0
	github.com/tinylib/msgp v1.6.4
)

require (
	github.com/coder/websocket v1.8.14 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
)
