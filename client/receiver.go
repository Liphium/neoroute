package client

type Receiver struct {
	Sender
}

// NewReceiver returns and initialized receiver.
// Use this when you not only want to send, but receive too.
func NewReceiver(config Config) *Receiver {
	return &Receiver{
		Sender: NewSender(config),
	}
}

func (r *Receiver) setEvent(eventName string, receiveFunc func(c *Ctx)) {
	r.setReceiver(eventName, receiveFunc)
}
