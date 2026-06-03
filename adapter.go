package neoroute

type Adapter interface {
	send(b []byte) error
	setRemoveFunc(removeFunc func())
}
