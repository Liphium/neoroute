package tui

import (
	"charm.land/bubbles/v2/key"
)

type FullKeyHelp struct {
	Title string
	Keys  [][]key.Binding
}

type keyProvider interface {
	Children() []keyProvider

	FooterKeys() []key.Binding
	FullKeyHelp() FullKeyHelp
}
