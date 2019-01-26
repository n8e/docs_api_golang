package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/n8e/rest-api/handlers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/people", handlers.PersonsHandler).Methods("GET")
	router.HandleFunc("/people/{id}", handlers.PersonHandler).Methods("GET")
	router.HandleFunc("/people/{id}", handlers.CreatePersonHandler).Methods("POST")
	router.HandleFunc("/people/{id}", handlers.DeletePersonHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}
