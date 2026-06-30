module web_socket

go 1.26.4

replace github.com/Liphium/neoroute => ../../

replace github.com/Liphium/neoroute/client => ../../client/

replace github.com/Liphium/neoroute/client/transporter/websocket => ../../client/transporter/websocket

replace github.com/Liphium/neoroute/transporter/websocket => ../../transporter/websocket

require (
	github.com/Liphium/neoroute v0.2.0
	github.com/coder/websocket v1.8.15 // indirect
	github.com/tinylib/msgp v1.6.4
)

require (
	github.com/Liphium/neoroute/client v0.0.0
	github.com/Liphium/neoroute/client/transporter/websocket v0.0.0-00010101000000-000000000000
	github.com/Liphium/neoroute/transporter/websocket v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
)
