module web_transport_server

go 1.26

replace github.com/Liphium/neoroute => ../../

replace github.com/Liphium/neoroute/transporter/web_transport => ../../transporter/web_transport

require github.com/quic-go/quic-go v0.60.0

require (
	github.com/Liphium/neoroute v0.0.0
	github.com/Liphium/neoroute/transporter/web_transport v0.0.0-00010101000000-000000000000
	github.com/quic-go/webtransport-go v0.11.0
	github.com/tinylib/msgp v1.6.4
)

require (
	github.com/dunglas/httpsfv v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
)
