package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-pg/pg/v10"
)

type Profile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {

	// Create a connection pool with a maximum of 10 connections
	opt := &pg.Options{
		Addr:     "localhost:5432",
		User:     "go",
		Password: "go",
		Database: "benchmark",
		PoolSize: 5, // Set the maximum number of connections to 10
	}
	db := pg.Connect(opt)
	defer db.Close()

	http.HandleFunc("/go/fetch", func(w http.ResponseWriter, r *http.Request) {
		fetchProfiles(w, r, db)
	})

	done := make(chan bool)

	go func() {
		// Start the HTTP server
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Trigger the callback once the server starts listening
	go func() {
		_, err := http.Get("http://localhost:8080") // Replace with your server URL
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Server started listening on port 8080")
		// Perform any additional actions or logic here
		done <- true
	}()

	<-done // Wait for the callback to complete

	// Keep the main goroutine alive
	select {}
}

func fetchProfiles(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	// Query all rows from the profile table
	var profiles []Profile
	_, err := db.Query(&profiles, "SELECT * FROM profile")
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
