package neoroute

type Config struct {
	ErrorHandler func(err error) string
}
