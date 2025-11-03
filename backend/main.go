package main

import (
	"encoding/json" // for encoding/decoding JSON
	"fmt"
	"io"       // for reading from HTTP responses
	"log"      // for logging errors/info
	"net/http" // for HTTP server and client
	"os"
	"sync" // for mutex locks (protects cache from concurrent HTTP requests)
	"time" // for cache timestamps

	"github.com/joho/godotenv"
)

// Cache for station data (SERVER-SIDE)
// This cache is shared by ALL users, when one user triggers a cache refresh, everyone benefits.
// Only fetches from WMATA API once every 24 hours.
var (
	cachedStations  []StationInfo     // Cached station data (stored in server memory)
	cachedEntrances []StationEntrance // Cached entrance data (stored in server memory)
	cacheTime       time.Time         // When the cache was last updated
	cacheDuration   = 24 * time.Hour  // Cache for 24 hours (station data rarely changes)
	cacheMutex      sync.RWMutex      // Protects cache from concurrent HTTP requests
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

// StationEntrance struct: Info about a station entrance
type StationEntrance struct {
	Description  string  `json:"Description"`
	ID           string  `json:"ID"`
	Lat          float64 `json:"Lat"`
	Lon          float64 `json:"Lon"`
	Name         string  `json:"Name"`
	StationCode1 string  `json:"StationCode1"`
	StationCode2 string  `json:"StationCode2"`
}

// EntrancesResponse struct: holds all station entrances
type EntrancesResponse struct {
	Entrances []StationEntrance `json:"Entrances"`
}

// Helper function to handle CORS and preflight requests
// CORS = Cross-Origin Resource Sharing. Browsers block requests between different origins (different ports/domains) for security.
// Our frontend runs on localhost:3000, backend on localhost:8080, different origins.
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

// Helper function to fetch all stations with caching
// Returns cached data if it's fresh, otherwise fetches from API
func fetchAllStations(apiKey string) ([]StationInfo, error) {
	// Check if cache is still valid (using read lock for concurrent safety)
	cacheMutex.RLock()
	if time.Since(cacheTime) < cacheDuration && len(cachedStations) > 0 {
		log.Println("Returning cached station data")
		defer cacheMutex.RUnlock()
		return cachedStations, nil
	}
	cacheMutex.RUnlock()

	// Cache is stale or empty, fetch fresh data
	log.Println("Cache expired or empty, fetching fresh station data from API")

	// Acquire write lock to update cache
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Double-check: another goroutine might have updated cache while we waited for the lock
	if time.Since(cacheTime) < cacheDuration && len(cachedStations) > 0 {
		return cachedStations, nil
	}

	// Fetch basic station list first
	body, err := fetchFromWMATA("https://api.wmata.com/Rail.svc/json/jStations", apiKey)
	if err != nil {
		return nil, err
	}

	var stationsResp StationsResponse
	err = json.Unmarshal(body, &stationsResp)
	if err != nil {
		return nil, err
	}

	// Fetch detailed info for each station sequentially (avoids rate limiting)
	var detailedStations []StationInfo
	startTime := time.Now()
	log.Printf("Fetching details for %d stations...\n", len(stationsResp.Stations))

	for _, station := range stationsResp.Stations {
		url := fmt.Sprintf("https://api.wmata.com/Rail.svc/json/jStationInfo?StationCode=%s", station.Code)
		infoBody, err := fetchFromWMATA(url, apiKey)
		if err != nil {
			log.Printf("Error fetching info for %s: %v\n", station.Code, err)
			continue
		}

		var stationInfo StationInfo
		err = json.Unmarshal(infoBody, &stationInfo)
		if err != nil {
			log.Printf("Error parsing info for %s: %v\n", station.Code, err)
			continue
		}

		detailedStations = append(detailedStations, stationInfo)
	}

	totalDuration := time.Since(startTime)
	log.Printf("Successfully fetched %d stations in %v (avg: %v per station)\n",
		len(detailedStations), totalDuration, totalDuration/time.Duration(len(detailedStations)))

	// Fetch all station entrances (no parameters = get all)
	log.Println("Fetching station entrances...")
	entrancesBody, err := fetchFromWMATA("https://api.wmata.com/Rail.svc/json/jStationEntrances", apiKey)
	if err != nil {
		log.Println("Error fetching entrances:", err)
		// Don't fail completely if entrances fail, just log it
	} else {
		var entrancesResp EntrancesResponse
		err = json.Unmarshal(entrancesBody, &entrancesResp)
		if err != nil {
			log.Println("Error parsing entrances:", err)
		} else {
			cachedEntrances = entrancesResp.Entrances
			log.Printf("Cached %d station entrances\n", len(cachedEntrances))
		}
	}

	// Update cache
	cachedStations = detailedStations
	cacheTime = time.Now()
	log.Printf("Cached %d stations at %s\n", len(cachedStations), cacheTime.Format(time.RFC3339))

	return detailedStations, nil
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("WMATA_API_KEY") // Get API key

	fmt.Println("Server running on :8080")

	// Pre-warm the cache on server startup
	// Runs in background (goroutine) so server starts immediately
	// This way the first user gets instant results instead of waiting 30-60 seconds
	log.Println("Pre-warming cache with station data...")
	go func() {
		_, err := fetchAllStations(apiKey)
		if err != nil {
			log.Println("Failed to pre-warm cache:", err)
		} else {
			log.Println("Cache pre-warmed successfully!")
		}
	}()

	// Handler for /stations - returns ALL station details with coordinates
	// Uses server-side cache: instant for all users as long as cache is fresh (24 hours)
	http.HandleFunc("/stations", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Use the caching function instead of fetching every time
		detailedStations, err := fetchAllStations(apiKey)
		if err != nil {
			log.Println("Error fetching stations:", err)
			http.Error(w, "API fetch failed", 500)
			return
		}

		// Return all detailed station info in one response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(detailedStations)
	})

	// Handler for /entrances - returns station entrances for a specific station code
	http.HandleFunc("/entrances", func(w http.ResponseWriter, r *http.Request) {
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

		// Ensure cache is populated (will use cached data if available)
		_, err := fetchAllStations(apiKey)
		if err != nil {
			log.Println("Error ensuring cache:", err)
			http.Error(w, "Cache fetch failed", 500)
			return
		}

		// Filter entrances for this station code
		cacheMutex.RLock()
		var stationEntrances []StationEntrance
		for _, entrance := range cachedEntrances {
			if entrance.StationCode1 == stationCode || entrance.StationCode2 == stationCode {
				stationEntrances = append(stationEntrances, entrance)
			}
		}
		cacheMutex.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stationEntrances)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
