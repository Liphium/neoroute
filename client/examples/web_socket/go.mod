module web_socket

go 1.26.3

replace github.com/Liphium/neoroute/client => ../../

require (
	github.com/tinylib/msgp v1.6.4
	github.com/Liphium/neoroute/client v0.0.0
)

require github.com/gorilla/websocket v1.5.3 // indirect

require (
	github.com/philhofer/fwd v1.2.0 // indirect
)
