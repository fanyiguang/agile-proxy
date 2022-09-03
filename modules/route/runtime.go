package route

import jsoniter "github.com/json-iterator/go"

var (
	route = make(map[string]Route)
	json  = jsoniter.ConfigCompatibleWithStandardLibrary
)

func GetRoute(name string) (t Route) {
	return route[name]
}

func GetAllRoute() map[string]Route {
	return route
}

func CloseAllRoutes() {
	for _, _route := range route {
		if _route != nil {
			_ = _route.Close()
		}
	}
}
