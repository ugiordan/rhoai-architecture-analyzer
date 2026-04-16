package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
)

func saveToDB(table string, data map[string]interface{}) error {
	db, _ := sql.Open("postgres", "host=localhost")
	_, err := db.Exec(fmt.Sprintf("INSERT INTO %s VALUES ($1)", table), data)
	return err
}

func getFromDB(table string) []map[string]interface{} {
	db, _ := sql.Open("postgres", "host=localhost")
	rows, _ := db.Query(fmt.Sprintf("SELECT * FROM %s", table))
	_ = rows
	return nil
}

func callExternalAPI(url string, payload []byte) error {
	_, err := http.Post(url, "application/json", bytes.NewReader(payload))
	return err
}
