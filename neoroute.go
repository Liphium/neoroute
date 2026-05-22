package neoroute

func NewNeoRouter[D any](config Config) *NeoRouter[D] {
	return &NeoRouter[D]{
		routes:     make(map[string]func(c *ctx[D]) error),
		middleware: make(map[string]func(c *ctx[D]) bool),
		config:     config,
	}
}
