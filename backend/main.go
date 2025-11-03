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

// Address struct: holds address information for a station
type Address struct {
	City   string `json:"City"`
	State  string `json:"State"`
	Street string `json:"Street"`
	Zip    string `json:"Zip"`
}

// StationInfo struct: Detailed info about a single station
type StationInfo struct {
	Address          Address `json:"Address"`
	Code             string  `json:"Code"`
	Lat              float64 `json:"Lat"`
	Lon              float64 `json:"Lon"`
	LineCode1        string  `json:"LineCode1"`
	LineCode2        string  `json:"LineCode2"`
	LineCode3        string  `json:"LineCode3"`
	LineCode4        string  `json:"LineCode4"`
	Name             string  `json:"Name"`
	StationTogether1 string  `json:"StationTogether1"`
	StationTogether2 string  `json:"StationTogether2"`
}

// Helper function to handle CORS and preflight requests
// CORS = Cross-Origin Resource Sharing. Browsers block requests between different origins (different ports/domains) for security.
// Our frontend runs on localhost:3000, backend on localhost:8080 — different origins!
// These headers tell the browser: "Hey, it's okay, I allow localhost:3000 to access my data."
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Only allow our frontend (for production, use specific origin, not "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")         // Allow GET requests and OPTIONS (preflight)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")         // Allow Content-Type header
}

// Helper function to fetch from WMATA API
// Takes a URL and API key, returns the response body as bytes or an error
func fetchFromWMATA(url string, apiKey string) ([]byte, error) {
	// Build a GET request to the WMATA API
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("api_key", apiKey)

	// &http.Client{} means "create a new http.Client and give me a pointer to it"
	client := &http.Client{}
	resp, err := client.Do(req) // Send the request
	if err != nil {
		return nil, err
	}
	// defer = "run this when the function returns" (cleanup); always closes the response body, even if there's an error or early return
	defer resp.Body.Close()

	// Read the response body (JSON)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if the API returned a success status code (200 OK)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return body, nil
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("WMATA_API_KEY") // Get API key
	//url := "https://api.wmata.com/Rail.svc/json/jStations" // API endpoint

	// Handler for /stations
	http.HandleFunc("/stations", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		body, err := fetchFromWMATA("https://api.wmata.com/Rail.svc/json/jStations", apiKey)
		if err != nil {
			log.Println("Error fetching stations:", err)
			http.Error(w, "API fetch failed", 500)
			return
		}

		var stationsResp StationsResponse
		err = json.Unmarshal(body, &stationsResp)
		if err != nil {
			log.Println("JSON unmarshal error:", err)
			http.Error(w, "JSON parse failed", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stationsResp.Stations)
	})

	// Handler for /station-info
	http.HandleFunc("/station-info", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		stationCode := r.URL.Query().Get("code")
		if stationCode == "" {
			http.Error(w, "Missing station code", 400)
			return
		}

		url := fmt.Sprintf("https://api.wmata.com/Rail.svc/json/jStationInfo?StationCode=%s", stationCode)
		body, err := fetchFromWMATA(url, apiKey)
		if err != nil {
			log.Println("Error fetching station info:", err)
			http.Error(w, "API fetch failed", 500)
			return
		}

		var stationInfo StationInfo
		err = json.Unmarshal(body, &stationInfo)
		if err != nil {
			log.Println("JSON unmarshal error:", err)
			http.Error(w, "JSON parse failed", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stationInfo)
	})

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
