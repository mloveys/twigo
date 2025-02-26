package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	expectedToken string
	gotifyToken   string
	gotifyURL     string
)

// GotifyMessage defines the JSON structure that Gotify expects. TODO: Add priority
type GotifyMessage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Extract token from query parameters.
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Println("Error: Missing token in request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify the token.
	if token != expectedToken {
		log.Println("Error: Unauthorized token attempt")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Only accept POST requests.
	if r.Method != http.MethodPost {
		log.Printf("Error: Method %s not allowed", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse form data from the POST body.
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form data: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve form values.
	from := r.FormValue("From")
	bodyMsg := r.FormValue("Body")
	if from == "" || bodyMsg == "" {
		log.Println("Error: Missing required form parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Build the Gotify message.
	gotifyMsg := GotifyMessage{
		Title:   "Message from " + from,
		Message: bodyMsg,
	}

	// Marshal the Gotify message to JSON.
	msgBytes, err := json.Marshal(gotifyMsg)
	if err != nil {
		log.Printf("Error creating message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the Gotify URL construction
	url := fmt.Sprintf("%s/message?token=%s", gotifyURL, gotifyToken)

	// Make the POST request to the Gotify endpoint.
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		log.Printf("Failed to send message to Gotify: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read and log the response from Gotify
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading Gotify response: %v", err)
	} else {
		log.Printf("Gotify response: %s", respBody)
	}

	// Return only the status code from Gotify
	w.WriteHeader(resp.StatusCode)
}

func main() {
	// Get tokens from environment variables
	expectedToken = os.Getenv("WEBHOOK_TOKEN")
	if expectedToken == "" {
		log.Fatal("WEBHOOK_TOKEN environment variable not set")
	}

	gotifyToken = os.Getenv("GOTIFY_TOKEN")
	if gotifyToken == "" {
		log.Fatal("GOTIFY_TOKEN environment variable not set")
	}

	gotifyURL = os.Getenv("GOTIFY_URL")
	if gotifyURL == "" {
		log.Fatal("GOTIFY_URL environment variable not set")
	}

	http.HandleFunc("/receive/", handler)
	log.Println("Server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

