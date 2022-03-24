package routes

import (
	"api/src/controllers"
	"net/http"
)

var clusterRoutes = []Route{

	{
		URI:      "/cluster",
		Method:   http.MethodGet,
		Function: controllers.GetClusterInfo,
		auth:     false,
	},
}
