package routes

import (
	"net/http"
)

type NamedRoute struct {
	Name        string
	Pattern     string
	StripPrefix string
	Handler     func(http.ResponseWriter, *http.Request)
}

func SetupRoutes(
	namedRoutes ...*NamedRoute,
) *http.ServeMux {
	mux := http.NewServeMux()

	for i := 0; i < len(namedRoutes); i++ {
		route := namedRoutes[i]
		mux.HandleFunc(route.Pattern, route.Handler)
	}

	return mux
}
