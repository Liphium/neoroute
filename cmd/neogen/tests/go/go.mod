module neogen-runtime-tests

go 1.27

require (
	github.com/Liphium/neoroute/client/sdk/v2 v2.0.0
	github.com/Liphium/neoroute/client/transporter/http/v2 v2.0.0-00010101000000-000000000000
	github.com/Liphium/neoroute/client/transporter/websocket/v2 v2.0.0-00010101000000-000000000000
	github.com/Liphium/neoroute/v2 v2.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.12.1
	github.com/tinylib/msgp v1.6.4
)

require (
	github.com/coder/websocket v1.8.15 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)

replace github.com/Liphium/neoroute/v2 => ../../../..

replace github.com/Liphium/neoroute/client/sdk/v2 => ../../../../client/sdk

replace github.com/Liphium/neoroute/client/transporter/http/v2 => ../../../../client/transporter/http

replace github.com/Liphium/neoroute/client/transporter/websocket/v2 => ../../../../client/transporter/websocket
