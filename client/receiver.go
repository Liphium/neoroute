package client

type Receiver struct {
	Handler
}

func NewReceiver(config Config) *Receiver {
	return &Receiver{
		Handler: *NewHandler(config),
	}
}
