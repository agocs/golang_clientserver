package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/agocs/golang_clientserver/payload"
)

func main() {
	// Create a new HTTP server
	server := &http.Server{
		Addr:    ":8080",                   // Set the address to listen on
		Handler: http.HandlerFunc(handler), // Set the handler function
	}

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		panic(err) // Handle any errors that occur while starting the server
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // Record the start time
	// Set the content type to JSON
	log.Printf("request started at second: %d, ns: %d", start.Second(), start.Nanosecond())
	log.Printf("\n\n\n")

	reqPayload := payload.Payload{}
	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		log.Printf("error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	payloadTime := time.Now()
	log.Printf("payload sent at %s", reqPayload.SentTime.Format(time.RFC3339))
	log.Printf("payload received at second: %d, ns: %d", payloadTime.Second(), payloadTime.Nanosecond())
	log.Printf("duration to decode payload: %s", payloadTime.Sub(start).String())
	log.Printf("\n\n\n")

	w.Header().Set("Content-Type", "application/json")

	// Write a simple JSON response
	w.Write([]byte(`{"message": "Hello, World!"}`))
	log.Printf("request completed in %s", time.Since(start).String())
	log.Printf("\n\n\n")
}
