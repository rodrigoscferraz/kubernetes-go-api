package routes

import (
	"net/http"
	"orion-api/src/controllers"

	"github.com/gorilla/mux"
)

type Route struct {
	URI      string
	Method   string
	Function func(http.ResponseWriter, *http.Request)
	auth     bool
}


// Configure put all routes inside the router
func Configure(r *mux.Router, kubeClient controllers.KubeClient) *mux.Router {
	
	var routes = []Route{
		{
			URI:      "/cluster",
			Method:   http.MethodGet,
			Function: kubeClient.GetClusterInfo,
			auth:     false,
		},
	}

	for _, route := range routes {
		r.HandleFunc(route.URI, route.Function).
			Methods(route.Method)
	}

	return r
}
