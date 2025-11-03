package main

import (
	"encoding/json" // for encoding/decoding JSON
	"fmt"
	"io"       // for reading from HTTP responses
	"log"      // for logging errors/info
	"net/http" // for HTTP server and client
	"os"

	"github.com/joho/godotenv"
)

// Station struct: like a struct in Rust, defines fields and their types
type Station struct {
	Name string `json:"Name"` // field maps to "Name" in JSON
	Code string `json:"Code"` // field maps to "Code" in JSON
}

// StationsResponse struct: holds a slice (dynamic array, like Vec in Rust) of Station
type StationsResponse struct {
	Stations []Station `json:"Stations"` // maps to "Stations" in JSON
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("WMATA_API_KEY")                   // Get API key
	url := "https://api.wmata.com/Rail.svc/json/jStations" // API endpoint

	// Set up a handler for the /stations route.
	// When someone visits /stations, this function runs.
	http.HandleFunc("/stations", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers to allow requests from the frontend
		// CORS = Cross-Origin Resource Sharing. Browsers block requests between different origins (different ports/domains) for security.
		// Our frontend runs on localhost:3000, backend on localhost:8080 — different origins!
		// These headers tell the browser: "Hey, it's okay, I allow localhost:3000 to access my data."
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Only allow our frontend (for production, use specific origin, not "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")         // Allow GET requests and OPTIONS (preflight)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")         // Allow Content-Type header

		// Handle preflight OPTIONS request
		// Before the real request, browsers sometimes send a "preflight" OPTIONS request to ask: "Is this allowed?"
		// We just respond with 200 OK to say "yes, it's allowed."
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// The "w" in this function is the response writer.
		// We use "w" to write data back to client (whoever made the HTTP request, like our frontend)
		//
		// The "r" is the incoming HTTP request.
		// "r" contains all the information about what the client sent to us:
		//   - The URL and path they requested
		//   - Any query parameters (like /stations?foo=bar)
		//   - HTTP headers (like User-Agent, etc.)
		//   - The HTTP method (GET, POST, etc.)
		//   - Any data they sent in the request body (for POST/PUT)
		// You use "r" if you want to read what the client sent, like checking query params or headers.

		// In Go, most functions that can fail return an error as the last return value.
		// You must check for errors after every operation that can fail.
		// This is different from languages that use exceptions or the ? operator (like Rust).
		// If you don't check for errors, your program might continue with bad data or crash later.
		// By handling errors right away, you can log them and send a proper response to the client.

		// Build a GET request to the WMATA API.
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			// log.Fatal prints the error and exits the program
			log.Fatal(err)
		}
		req.Header.Set("api_key", apiKey) // Add API key to the request headers

		client := &http.Client{}
		resp, err := client.Do(req) // Send the request
		if err != nil {
			// log the error and send an HTTP 500 error to the client, remember, client will be web dashboard!
			log.Println("Error fetching API:", err)
			http.Error(w, "API fetch failed", 500)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body) // Read the response body (JSON)
		if err != nil {
			log.Println("Error reading body:", err)
			http.Error(w, "Read failed", 500)
			return
		}

		// Check if the API returned a success status code (200 OK)
		if resp.StatusCode != 200 {
			log.Println("API returned non-200:", resp.StatusCode)
			http.Error(w, "API error", 500)
			return
		}

		// Parse JSON response
		// We declare a variable of type StationsResponse (our struct that matches the JSON structure)
		var stationsResp StationsResponse

		// json.Unmarshal takes the raw JSON bytes (body) and fills our struct with the data.
		// This is like serde_json::from_str in Rust.
		// The fields in the JSON must match the struct fields (including tags like `json:"Stations"`).
		// After this, stationsResp will have all the station data from the API.
		err = json.Unmarshal(body, &stationsResp)
		if err != nil {
			log.Println("JSON unmarshal error:", err)
			http.Error(w, "JSON parse failed", 500)
			return
		}

		// Now stationsResp is a Go struct with all the data from the JSON.
		// You can use it in Go code, or send it back to the frontend as JSON again.

		// Return JSON to client
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stationsResp.Stations)
	})

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	// This starts the HTTP server, if there's an error, log.Fatal will print the error and exit the program
}

// ---
// What's the difference between fmt and log?
//
// - fmt is for printing stuff to your own terminal/console, like "hey, server started!".
//   It doesn't add timestamps or anything fancy, just plain output for you.
//   Example: fmt.Println("Hello, world!")
//
// - log is for logging errors or important info, especially for debugging or when something goes wrong.
//   log adds timestamps and is meant for tracking issues, not just printing random info.
//   log.Fatal is like panic! -- it prints the error and then stops the program.
//   Example: log.Println("Something went wrong:", err)
//
// TL;DR: use fmt for your own info, use log for errors and things you want to track or debug.
