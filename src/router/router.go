package router

import (
	"orion-api/src/controllers"
	"orion-api/src/router/routes"

	"github.com/gorilla/mux"
)

func Generate(kubeClient controllers.KubeClient) *mux.Router {
	r := mux.NewRouter()
	return routes.Configure(r, kubeClient)
}
