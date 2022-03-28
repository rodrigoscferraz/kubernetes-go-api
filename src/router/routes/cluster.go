package routes

import (
	"net/http"
	"orion-api/src/controllers"
)

var clusterRoutes = []Route{

	{
		URI:      "/cluster",
		Method:   http.MethodGet,
		Function: controllers.GetClusterInfo,
		auth:     false,
	},
}
