module http_server

go 1.26.4

replace github.com/Liphium/neoroute => ../../

replace github.com/Liphium/neoroute/transporter/http => ../../transporter/http

require (
	github.com/Liphium/neoroute v0.2.0
	github.com/Liphium/neoroute/transporter/http v0.0.0-00010101000000-000000000000
	github.com/tinylib/msgp v1.6.4
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
)
