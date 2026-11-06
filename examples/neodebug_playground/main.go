// Neodebug editor playground. The HTTP transporter is never served.
package main

import (
	"net/http"

	"github.com/Liphium/neoroute/pkg/neodebug/v2"
	"github.com/Liphium/neoroute/pkg/neodebug/v2/config"
	http_transporter "github.com/Liphium/neoroute/transporter/http/v2"
	"github.com/Liphium/neoroute/v2"
	"github.com/Liphium/neoroute/v2/neoschema"
)

//go:generate msgp -file models.go -tests=false

func main() {
	g := newGenerator()
	// Must run before Neodebug: its schema subprocess invokes this same program.
	g.PrintIfDesired()
	neodebug.Run(config.DebugConfig{
		TransporterName: "http",
		TransporterURL:  "http://127.0.0.1:6122/http",
	})
}

func newGenerator() *neoschema.Generator {
	router := neoroute.NewRouter(neoroute.Config[neoroute.NoData]{})
	router.RouteOk("scalars", Scalars)
	router.RouteOk("collections", Collections)
	router.RouteOk("nullable", Nullable)
	router.RouteOk("all-types", AllTypes)

	// Discard the HTTP hook: this transporter exists solely for schema generation.
	_, transporter := http_transporter.NewTransporter(router, http_transporter.Config[neoroute.NoData]{
		HandshakeFunc: func(*http.Request) (neoroute.NoData, bool) {
			return neoroute.NoData{}, true
		},
	})
	g := neoschema.NewGenerator()
	g.Transporter("http", transporter)
	return g
}

func Scalars(c *neoroute.OkCtx[neoroute.NoData], _ ScalarRequest) error {
	return c.RespondOk()
}

func Collections(c *neoroute.OkCtx[neoroute.NoData], _ CollectionRequest) error {
	return c.RespondOk()
}

func Nullable(c *neoroute.OkCtx[neoroute.NoData], _ NullableRequest) error {
	return c.RespondOk()
}

func AllTypes(c *neoroute.OkCtx[neoroute.NoData], _ AllTypesRequest) error {
	return c.RespondOk()
}
