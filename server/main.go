package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/agocs/golang_clientserver/payload"
)

func main() {
	// Create a new HTTP server
	// start a pprof server

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	server := &http.Server{
		Addr:    ":8080",                   // Set the address to listen on
		Handler: http.HandlerFunc(handler), // Set the handler function
	}

	log.Printf("Server starting at %s", time.Now().Format(time.RFC3339))
	// Start the server
	if err := server.ListenAndServe(); err != nil {
		panic(err) // Handle any errors that occur while starting the server
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // Record the start time
	log.Printf("request started at 			%s", start.Format("15:04:05.000000000"))

	reqPayload := payload.Payload{}
	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		log.Printf("error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	payloadTime := time.Now()
	log.Printf("payload sent at 			%s", reqPayload.SentTime.Format("15:04:05.000000000"))
	log.Printf("payload fully received at 		%s", payloadTime.Format("15:04:05.000000000"))
	log.Printf("duration to decode payload: 	%s", payloadTime.Sub(start).String())
	fmt.Print("\n\n\n")

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{"message": "Hello, World!"}`))
	log.Printf("request completed in %s", time.Since(start).String())
	log.Printf("\n\n\n")
}
