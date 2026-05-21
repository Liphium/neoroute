package neoroute

//go:generate msgp
type Request struct {
	// TODO: add other request fields
	Data []byte `msg:"data"`
}
