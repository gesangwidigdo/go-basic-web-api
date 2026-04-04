package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func resJson(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Set("Content-Type", "application/json")
	
	response := map[string]string{
		"message": message,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resJson(w, r, "welcome")
	})

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		resJson(w, r, "Hello World")
	})

	http.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len("/hello/"):]

		nameString := fmt.Sprintf("Hello %s", name)

		resJson(w, r, nameString)
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		resJson(w, r, "pong")
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}