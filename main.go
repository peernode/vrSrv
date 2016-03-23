package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func main() {
	router := httprouter.New()
	initPathRouter(router)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  1000 * time.Second,
		WriteTimeout: 500 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
