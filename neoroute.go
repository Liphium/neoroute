package neoroute

func NewNeoRouter(config Config) *NeoRouter {
	return &NeoRouter{
		routes:     make(map[string]func(c *ctx) error),
		middleware: make(map[string]func(c *ctx) bool),
		config:     config,
	}
}
