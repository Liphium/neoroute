package main

import (
	"log/slog"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/Liphium/neoroute"
	"github.com/Liphium/neoroute/neoschema"
	http_transporter "github.com/Liphium/neoroute/transporter/http"
	websocket_transporter "github.com/Liphium/neoroute/transporter/websocket"
	"github.com/coder/websocket"

	"github.com/Liphium/neoroute/pkg/neodebug"
	"github.com/Liphium/neoroute/pkg/neodebug/config"
)

//go:generate msgp

type SendRequest struct {
	Text   string `msg:"text"`
	Sender string `msg:"sender"`
}

type MessageEvent struct {
	Text      string `msg:"text"`
	Sender    string `msg:"sender"`
	Timestamp int64  `msg:"timestamp"`
}

// Create event registry and register event for message broadcast
var eventRegistry = neoroute.NewEventRegistry()
var createMessageEvent = eventRegistry.Register[MessageEvent]("message")

var adapterRegistry = neoroute.NewAdapterRegistry()

func main() {

	if slices.Contains(os.Args, "--debug-http") {
		neodebug.Run(config.DebugConfig{

			// The name of the transporter you want to connect to (from your schema).
			TransporterName: "http",

			// The URL to that transporter.
			TransporterURL: "http://localhost:6121/http",
		})
	}

	if slices.Contains(os.Args, "--debug-ws") {
		neodebug.Run(config.DebugConfig{

			// The name of the transporter you want to connect to (from your schema).
			TransporterName: "ws",

			// The URL to that transporter.
			TransporterURL: "ws://localhost:6121/ws",
		})
	}

	// Create shared router for http and websocket transporter
	router := neoroute.NewRouter(neoroute.Config[neoroute.NoData]{
		ErrorHandler: func(err error, c *neoroute.Ctx[neoroute.NoData]) string {
			slog.Info("Error in neoroute", "error", err)
			return "Something went wrong."
		},
	})

	// Register route for sending messages
	router.RouteOk("send", Send)

	// Initialize the transporters after the routes have been registered.
	mux := InitTransporter(router)

	slog.Info("Starting neoroute chat server example")

	slog.Info("listening on localhost:6121")
	if err := http.ListenAndServe(":6121", mux); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func Send(c *neoroute.OkCtx[neoroute.NoData], req SendRequest) error {
	if req.Text == "" {
		return neoroute.NewError("text is required")
	}

	slog.Info("Received new message", "text", req.Text, "sender", req.Sender)

	// Broadcast the message to all connected clients
	adapterRegistry.Broadcast(createMessageEvent(MessageEvent{
		Text:      req.Text,
		Sender:    req.Sender,
		Timestamp: time.Now().Unix(),
	}))

	return c.RespondOk()
}

func InitTransporter(router *neoroute.Router[neoroute.NoData]) *http.ServeMux {

	// Create HTTP transporter
	httpHook, httpT := http_transporter.NewTransporter(router, http_transporter.Config[neoroute.NoData]{
		HandshakeFunc: func(r *http.Request) (neoroute.NoData, bool) {
			return neoroute.NoData{}, false
		},
	})

	// Create WebSocket transporter
	wsHook, wsT := websocket_transporter.NewTransporter(router, websocket_transporter.Config[neoroute.NoData]{
		HandshakeFunc: func(r *http.Request) (neoroute.NoData, bool) {
			return neoroute.NoData{}, false
		},
		AcceptOptions: &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		},
		EnterNetworkFunc: func(session *neoroute.Session[neoroute.NoData]) {
			adapter, err := session.Adapt()
			if err != nil {
				slog.Info("Failed to adapt session", "error", err)
				session.Disconnect()
				return
			}
			adapterRegistry.Register(session.Id(), adapter)
		},
	})

	// Register transporters for schema generation
	g := neoschema.NewGenerator()
	g.Transporter("ws", wsT)
	g.Transporter("http", httpT)

	// Then let the program panic and print the schema when --neo-generate is set
	g.PrintIfDesired()

	// Register events with the WebSocket transporter
	wsT.AddEventRegistry(eventRegistry)

	// Mount HTTP and WebSocket transporter
	mux := http.NewServeMux()
	mux.HandleFunc("/http", httpHook)
	mux.HandleFunc("/ws", wsHook)

	return mux
}
