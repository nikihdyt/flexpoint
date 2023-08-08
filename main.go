package main

import (
	"flexpoint/connection"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	connection.ConnectDatabase()

	r := mux.NewRouter()

	log.Fatal(http.ListenAndServe(":9990", r))
}
