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

func applyRouterType(router Router[NoData], routerType string) Router[D] {
	switch routerType {
	case routerTypeNeo:
		return router.AddRouters(NewRouter(Config[D]{}))
	case routerTypeGroup:

	case routerTypeRoute:

	default:
		panic("unknown router type: " + routerType)
	}
}
