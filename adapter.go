package neoroute

type Adapter interface {
	isEventRegistered(name string) bool
	getTransportType() string
	send(b []byte) error
	setRemoveFunc(removeFunc func())
	disconnect()
}
