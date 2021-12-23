package main

import (
	"go-example/route"
	"log"
	"net/http"

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
