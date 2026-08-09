package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserResponse struct {
	Service string `json:"service"`
	Message string `json:"message"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", usersHandler)

	server := &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	log.Println("Upstream service running on port 4000")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	response := UserResponse{
		Service: "user-service",
		Message: "Hello from upstream",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
