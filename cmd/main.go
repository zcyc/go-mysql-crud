package main

import (
	"log"
	"net/http"

	"go-example/route"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	route.AddHandler(r)
	err := http.ListenAndServe(":80", r)
	if err != nil {
		log.Fatal(err)
	}
}
