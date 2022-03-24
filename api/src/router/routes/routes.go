package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI      string
	Method   string
	Function func(http.ResponseWriter, *http.Request)
	auth     bool
}

// Configure put all routes inside the router
func Configure(r *mux.Router) *mux.Router {
	routes := clusterRoutes

	for _, route := range routes {
		r.HandleFunc(route.URI, route.Function).
			Methods(route.Method)
	}

	return r
}
