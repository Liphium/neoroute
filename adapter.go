package neoroute

type Adapter interface {
	IsEventRegistered(name string) bool // DON'T USE THIS ONLY IMPLEMENT IT IF YOU ARE CREATING AN ADAPTER
	GetTransportType() string           // DON'T USE THIS ONLY IMPLEMENT IT IF YOU ARE CREATING AN ADAPTER
	Send(b []byte) error                // DON'T USE THIS ONLY IMPLEMENT IT IF YOU ARE CREATING AN ADAPTER
	SetRemoveFunc(removeFunc func())    // DON'T USE THIS ONLY IMPLEMENT IT IF YOU ARE CREATING AN ADAPTER
	Disconnect()                        // DON'T USE THIS ONLY IMPLEMENT IT IF YOU ARE CREATING AN ADAPTER
}
