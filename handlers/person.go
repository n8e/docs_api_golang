package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/n8e/rest-api/person"
)

var people []person.Person

func init() {
	people = append(people, person.Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &person.Address{City: "City X", State: "State X"}})
	people = append(people, person.Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &person.Address{City: "City Z", State: "State Y"}})
}

func PersonHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&person.Person{})
}

func PersonsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func CreatePersonHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person person.Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePersonHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}
