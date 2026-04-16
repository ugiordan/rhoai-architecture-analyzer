package main

import (
	"encoding/json"
	"net/http"
	"log"
)

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users := getFromDB("users")
	json.NewEncoder(w).Encode(users)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user map[string]interface{}
	json.NewDecoder(r.Body).Decode(&user)
	saveToDB("users", user)
	w.WriteHeader(http.StatusCreated)
}

func main() {
	http.HandleFunc("/users", handleGetUsers)
	http.HandleFunc("/users/create", handleCreateUser)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
