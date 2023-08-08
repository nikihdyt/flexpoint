package main

import (
	"flexpoint/connection"
	user_controller "flexpoint/controller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	connection.ConnectDatabase()

	r := mux.NewRouter()

	r.HandleFunc("/register", user_controller.RegisterUser).
		Methods("POST")

	r.HandleFunc("/users", user_controller.GetListUser).
		Methods("GET")
	r.HandleFunc("/user/{id}", user_controller.GetDetailUser).
		Methods("GET")

	r.HandleFunc("/event", user_controller.CreateEvent).
		Methods("POST")
	r.HandleFunc("/event/verify/{id}", user_controller.VerifyEventAndAddPoint).
		Methods("PUT")
	r.HandleFunc("/event/{id}", user_controller.GetEventDetail).
		Methods("GET")

	r.HandleFunc("/point", user_controller.GetListPoint).
		Methods("GET")
	r.HandleFunc("/reedem-voucher/{id}", user_controller.RedeemVoucher).
		Methods("POST")

	log.Fatal(http.ListenAndServe(":9990", r))
}
