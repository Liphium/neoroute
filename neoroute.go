package neoroute

func NewNeoRouter[D any](config Config) *NeoRouter[D] {
	return &NeoRouter[D]{
		routes:     make(map[string]func(c *Ctx[D]) error),
		middleware: make(map[string]func(c *Ctx[D]) bool),
		config:     config,
	}
}
