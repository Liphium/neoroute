package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Liphium/neoroute"
	http_transporter "github.com/Liphium/neoroute/transporter/http"
)

//go:generate msgp
type Request struct {
	Field1 string `msg:"field1"`
	Field2 int    `msg:"field2"`
}

//go:generate msgp
type Response struct {
	Field1 string `msg:"field1"`
	Field2 int    `msg:"field2"`
}

type SessionData struct {
	Token string
}

func main() {

	// Create router and set it for transporter
	router := neoroute.NewNeoRouter[SessionData](neoroute.Config{
		ErrorHandler: func(err error) string {
			return fmt.Sprintf("error: %v", err)
		},
	})

	// Create HTTP transporter
	hook, _ := http_transporter.NewHTTPTransporter(router, func(r *http.Request) (SessionData, bool) {

		// Set token if one provided as session data
		return SessionData{
			Token: r.URL.Query().Get("token"),
		}, true
	})

	router2 := neoroute.NewNeoRouter[SessionData](neoroute.Config{
		ErrorHandler: func(err error) string {
			return fmt.Sprintf("error: %v", err)
		},
	})

	// Create a router group, any function used with this router will be applied to all neo routes it contains.
	// This is useful when you want to apply certain groups, routes or middle ware to multiple routers at once.
	// You could use router2 for example for WebTransport and apply shared routes to both and add the WebTransport specific routes only to router2.
	rGroup := router.AddRouters(router2)

	// Route: simple.route
	// Wrap the RouteResponse call with a Use to directly apply a middleware to that route specifically.
	neoroute.Use(neoroute.RouteNoRequest(rGroup, "simple.route", func(c *neoroute.ResCtx[SessionData, Response, *Response]) error {
		return c.Respond(Response{Field1: "simple response that had no input", Field2: 68})
	}), "", func(c *neoroute.Ctx[SessionData]) bool {
		fmt.Println("middleware mounted directly on route was used")
		return true
	})

	// Create group for group1
	group1 := rGroup.Group("group1")

	// Apply auth middleware to group1
	// If the token provided in the handshake is not `secret_token` this wont let the user continue.
	// Now only simple.route can be accessed without a token.
	neoroute.Use(group1, "", func(c *neoroute.Ctx[SessionData]) bool {
		fmt.Printf("middleware for group1 used with route %v by userId %v\n", c.Route(), c.Session().Id())
		fmt.Printf("session data: %+v\n", c.Session().Data().Token)
		if c.Session().Data().Token == "" {
			fmt.Println("no token provided, rejecting request")
		}
		return c.Session().Data().Token == "secret_token"
	})

	// Create subroute for group1
	// Route: group1.route1
	neoroute.Route(group1, "route1", func(c *neoroute.ResCtx[SessionData, Response, *Response], req Request) error {
		return c.Respond(Response{
			Field1: "response to " + req.Field1,
			Field2: req.Field2 + 1,
		})
	})

	// Create group2 on top of group1
	group2 := group1.Group("group2")

	// Create subroute for group2
	// Route: group1.group2.route1
	neoroute.Route(group2, "route1", func(c *neoroute.ResCtx[SessionData, Response, *Response], req Request) error {
		return c.Respond(Response{
			Field1: "response to " + req.Field1,
			Field2: req.Field2 + 2,
		})
	})

	// Create server
	mux := http.NewServeMux()

	// Hook http transporter into / route
	mux.HandleFunc("/", hook)

	log.Println("listening on http://localhost:6121")
	if err := http.ListenAndServe(":6121", mux); err != nil {
		log.Fatal(err)
	}
}
