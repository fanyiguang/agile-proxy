package route

import jsoniter "github.com/json-iterator/go"

var (
	routes = make(map[string]Route)
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
)

func GetRoute(name string) (t Route) {
	return routes[name]
}

func GetAllRoute() map[string]Route {
	return routes
}

func CloseAllRoutes() {
	for _, _route := range routes {
		if _route != nil {
			_ = _route.Close()
		}
	}
}
