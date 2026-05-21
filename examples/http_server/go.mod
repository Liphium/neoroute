module http_server

go 1.26.3

replace github.com/Liphium/neoroute => ../../

require github.com/quic-go/quic-go v0.59.1

require (
	github.com/Liphium/neoroute v0.0.0
	github.com/tinylib/msgp v1.6.4
)

require (
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)
