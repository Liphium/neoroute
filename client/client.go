package client

type Ctx struct {
	data []byte // data field from Request struct
	name string // event name
}

func (c *Ctx) Data() []byte {
	return c.data
}

func (c *Ctx) Name() string {
	return c.name
}
