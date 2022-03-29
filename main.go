package main

import (
	"fmt"
	"log"
	"net/http"
	"orion-api/src/router"
)

func main() {
	fmt.Println("Running")

	r := router.Generate()

	log.Fatal(http.ListenAndServe(":8000", r))
}
