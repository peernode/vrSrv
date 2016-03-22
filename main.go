package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()
	initPathRouter(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
