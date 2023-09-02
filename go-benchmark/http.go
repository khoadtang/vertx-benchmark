package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Profile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	http.HandleFunc("/profiles", fetchProfiles)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchProfiles(w http.ResponseWriter, r *http.Request) {
	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", "postgres://go:go@localhost:5432/benchmark?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query all rows from the profile table
	rows, err := db.Query("SELECT id, name, age FROM profile")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create a slice to hold the profiles
	profiles := []Profile{}

	// Iterate over the rows and populate the profiles slice
	for rows.Next() {
		var profile Profile
		err := rows.Scan(&profile.ID, &profile.Name, &profile.Age)
		if err != nil {
			log.Fatal(err)
		}
		profiles = append(profiles, profile)
	}

	// Check for any errors during iteration
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// Convert profiles slice to JSON
	jsonData, err := json.Marshal(profiles)
	if err != nil {
		log.Fatal(err)
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response writer
	_, err = w.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
}
