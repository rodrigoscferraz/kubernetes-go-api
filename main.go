package main

import (
	"fmt"
	"log"
	"net/http"
	"orion-api/src/controllers"
	"orion-api/src/kube"
	"orion-api/src/router"
)

var port = "8000"
var kubeClient controllers.KubeClient

func init() {

	_, clientSet := kube.Kubeconf()

	kubeClient = controllers.KubeClient{
		Client: clientSet,
	}

}

func main() {
	fmt.Println("Running on port: ", port)

	r := router.Generate(kubeClient)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
