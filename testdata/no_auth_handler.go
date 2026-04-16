package main

import (
	"encoding/json"
	"net/http"
)

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	json.NewDecoder(r.Body).Decode(&req)
	deleteFromDB("users", req["id"].(string))
	w.WriteHeader(http.StatusOK)
}

func deleteFromDB(table string, id string) {
	// simulated DB delete
}
