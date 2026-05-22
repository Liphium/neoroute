package neoroute

type Transporter[D any] interface {
	SetRouter(r *NeoRouter[D])
}
