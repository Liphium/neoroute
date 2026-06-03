package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Liphium/neoroute"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Counter struct {
	mutex       *sync.Mutex
	echoCounter int
	puns        []string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin allows us to accept connections from different domains.
	// For development, we return true. For production, secure this!
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Setup events
var eventReg = neoroute.NewEventRegistry()
var CreateNewPunSubmittedEvent = neoroute.Register[NewPunEvent](eventReg, "new_pun_submitted")

func main() {
	counter := Counter{
		mutex: &sync.Mutex{},
	}

	adapterReg := neoroute.NewAdapterRegistry[struct{}]()

	hook, t := neoroute.NewWebSocketTransporter(neoroute.WSConfig[struct{}]{
		UpgradeFunc: upgrader.Upgrade,
		OverwriteSessionFunc: func(id string) bool {
			return true
		},
		HandshakeFunc: func(r *http.Request) (*neoroute.Session[struct{}], bool) {
			return neoroute.NewSession[struct{}](uuid.NewString()), true
		},
		EnterNetworkFunc: func(session *neoroute.Session[struct{}], t *neoroute.WebSocketTransporter[struct{}]) {

			log.Println("user connected")

			// Add to adapter registry, in this case we don't have to manually unregister the adapter, because we want then in the registry until they disconnect.
			// Then they will be removed automatically.
			adapter, err := t.Adapt(session.Id())
			if err != nil {
				log.Println("failed to create adapter for", session.Id(), "with error", err)
				return
			}
			adapterReg.Register(session.Id(), adapter)
		},
		DisconnectHandler: func(session *neoroute.Session[struct{}]) {
			log.Println("user disconnected")
		},
	})

	// Add events to transporter
	t.AddEventRegistry(eventReg)

	// Setup routes
	r := neoroute.NewNeoRouter[struct{}](neoroute.Config{
		ErrorHandler: func(err error) string {
			log.Println("error occurred: ", err)
			return "Internal server error."
		},
	})
	t.SetRouter(r)

	neoroute.Route(r, "echo", func(c *neoroute.ResCtx[EchoResponse, *EchoResponse, struct{}], req EchoRequest) error {
		log.Println("message received")
		counter.mutex.Lock()
		defer counter.mutex.Unlock()
		counter.echoCounter++
		return c.Respond(EchoResponse{
			RequestNumber: counter.echoCounter,
			Message:       req.Message,
		})
	})

	neoroute.RouteOk(r, "submit_pun", func(c *neoroute.OkCtx[struct{}], req SubmitPunRequest) error {
		log.Println("message received")
		// Check pun contains go
		if !strings.Contains(strings.ToLower(req.Pun), "go") {
			return c.RespondError(fmt.Errorf("pun has to contain at least one instance of go. For example How is it GOing."))
		}

		// Send new pun to all clients
		ev, err := CreateNewPunSubmittedEvent(NewPunEvent{
			Pun: req.Pun,
		})
		if err != nil {
			log.Println("failed to create pun event", err)
		}
		adapterReg.Broadcast(ev)

		counter.mutex.Lock()
		counter.puns = append(counter.puns, req.Pun)
		counter.mutex.Unlock()

		return c.RespondOk()
	})

	// Create websocket transporter and host it
	mux := http.NewServeMux()
	mux.HandleFunc("/", hook)

	log.Println("listening on localhost:6121")
	http.ListenAndServe(":6121", mux)
}
