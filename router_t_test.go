package neoroute

import "testing"

const (
	routerTypeNeo   = "neo"
	routerTypeGroup = "group"
	routerTypeRoute = "route"
)

type testRouter struct {
	router *Router[NoData]
}

func TestRouter(t *testing.T) {

}
