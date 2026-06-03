package client

import "time"

type Config struct {
	ErrorHandler   func(err error)
	RequestTimeout time.Duration
}
