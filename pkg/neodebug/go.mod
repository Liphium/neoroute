module github.com/Liphium/neoroute/pkg/neodebug/v2

go 1.27

replace github.com/Liphium/neoroute/v2 => ../../

replace github.com/Liphium/neoroute/client/sdk/v2 => ../../client/sdk

replace github.com/Liphium/neoroute/client/transporter/websocket/v2 => ../../client/transporter/websocket

replace github.com/Liphium/neoroute/client/transporter/http/v2 => ../../client/transporter/http

require (
	charm.land/bubbles/v2 v2.2.1
	charm.land/bubbletea/v2 v2.0.9
	charm.land/lipgloss/v2 v2.0.6
	github.com/Liphium/neoroute/client/sdk/v2 v2.0.0
	github.com/Liphium/neoroute/client/transporter/http/v2 v2.0.0
	github.com/Liphium/neoroute/client/transporter/websocket/v2 v2.0.0
	github.com/Liphium/neoroute/v2 v2.0.0
	github.com/Unbreathable/nowandlater v1.0.0
	github.com/stretchr/testify v1.12.1
	github.com/tinylib/msgp v1.6.4
)

require (
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/charmbracelet/colorprofile v0.4.3 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20260812204455-68fa937c71be // indirect
	github.com/charmbracelet/x/ansi v0.11.8 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/coder/websocket v1.8.15 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.1 // indirect
	github.com/mattn/go-runewidth v0.0.28 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/xo/terminfo v1.0.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)
