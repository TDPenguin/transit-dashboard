package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// Cache for station data (SERVER-SIDE)
// This cache is shared by ALL users, when one user triggers a cache refresh, everyone benefits.
// Only fetches from WMATA API once every 24 hours.
var (
	cachedStations  []StationInfo     // Cached station data (stored in server memory)
	cachedEntrances []StationEntrance // Cached entrance data (stored in server memory)
	cachedLines     []Lines           // Cached rail lines data
	cachedParking   []StationParking  // Cached parking data
	cacheTime       time.Time         // When the cache was last updated
	cacheDuration   = 24 * time.Hour  // Cache for 24 hours (station data rarely changes)
	cacheMutex      sync.RWMutex      // Protects cache from concurrent HTTP requests

	cachedPredictions       []TrainPrediction
	predictionCacheTime     time.Time
	predictionCacheDuration = 25 * time.Second // Cache valid for 25s (refreshed every 20s = 5s buffer)
	predictionMutex         sync.RWMutex
)

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

// Generic fetch and parse - combines fetch + unmarshal
func fetchAndParse(url string, apiKey string, target interface{}) error {
	body, err := fetchFromWMATA(url, apiKey)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// Helper function to fetch all stations with caching
// Returns cached data if it's fresh, otherwise fetches from API
func fetchAllStations(apiKey string) ([]StationInfo, error) {
	// Check if cache is still valid (using read lock for concurrent safety)
	cacheMutex.RLock()
	if time.Since(cacheTime) < cacheDuration && len(cachedStations) > 0 {
		defer cacheMutex.RUnlock()
		return cachedStations, nil
	}
	cacheMutex.RUnlock()

	return refreshAllStations(apiKey)
}

// refreshAllStations ALWAYS fetches fresh data (used by background refresh)
func refreshAllStations(apiKey string) ([]StationInfo, error) {
	fetchStart := time.Now()

	// Acquire write lock to update cache
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Double-check: someone might have just refreshed
	if time.Since(cacheTime) < 1*time.Minute && len(cachedStations) > 0 {
		return cachedStations, nil
	}

	// Fetch station list
	var stationsResp StationsResponse
	if err := fetchAndParse("https://api.wmata.com/Rail.svc/json/jStations", apiKey, &stationsResp); err != nil {
		return nil, err
	}

	// Fetch detailed info for each station (sequentially)
	var detailedStations []StationInfo
	for _, station := range stationsResp.Stations {
		url := fmt.Sprintf("https://api.wmata.com/Rail.svc/json/jStationInfo?StationCode=%s", station.Code)
		var stationInfo StationInfo
		if err := fetchAndParse(url, apiKey, &stationInfo); err != nil {
			log.Printf("ERROR fetching station %s: %v\n", station.Code, err)
			continue
		}
		detailedStations = append(detailedStations, stationInfo)
	}

	// Fetch station entrances
	var entrancesResp EntrancesResponse
	if err := fetchAndParse("https://api.wmata.com/Rail.svc/json/jStationEntrances", apiKey, &entrancesResp); err != nil {
		log.Printf("ERROR fetching entrances: %v\n", err)
	} else {
		cachedEntrances = entrancesResp.Entrances
	}

	// Fetch lines
	var linesResp LinesResponse
	if err := fetchAndParse("https://api.wmata.com/Rail.svc/json/jLines", apiKey, &linesResp); err != nil {
		log.Printf("ERROR fetching lines: %v\n", err)
	} else {
		cachedLines = linesResp.Lines
	}

	// Fetch parking
	var parkingResp StationsParkingResponse
	if err := fetchAndParse("https://api.wmata.com/Rail.svc/json/jStationParking", apiKey, &parkingResp); err != nil {
		log.Printf("ERROR fetching parking: %v\n", err)
	} else {
		cachedParking = parkingResp.StationsParking
	}

	// Update cache
	cachedStations = detailedStations
	cacheTime = time.Now()

	fetchDuration := time.Since(fetchStart)
	log.Printf("[Static] API calls: %dms, %d stations, %d entrances, %d lines, %d parking\n",
		fetchDuration.Milliseconds(), len(detailedStations), len(cachedEntrances), len(cachedLines), len(cachedParking))

	return detailedStations, nil
}

// Fetch train predictions with caching (20 second refresh)
func fetchTrainPredictions(apiKey string) ([]TrainPrediction, error) {
	predictionMutex.RLock()
	if time.Since(predictionCacheTime) < predictionCacheDuration && len(cachedPredictions) > 0 {
		defer predictionMutex.RUnlock()
		return cachedPredictions, nil
	}
	predictionMutex.RUnlock()

	return refreshTrainPredictions(apiKey)
}

// refreshTrainPredictions always fetches fresh data (used by background refresh)
func refreshTrainPredictions(apiKey string) ([]TrainPrediction, error) {
	predictionMutex.Lock()
	defer predictionMutex.Unlock()

	// Double-check pattern (someone might have just refreshed)
	if time.Since(predictionCacheTime) < 1*time.Second && len(cachedPredictions) > 0 {
		return cachedPredictions, nil
	}

	// Fetch fresh predictions
	fetchStart := time.Now()
	url := "http://api.wmata.com/StationPrediction.svc/json/GetPrediction/All"
	body, err := fetchFromWMATA(url, apiKey)
	if err != nil {
		return nil, err
	}
	fetchDuration := time.Since(fetchStart)

	var resp struct {
		Trains []TrainPrediction `json:"Trains"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	cachedPredictions = resp.Trains
	predictionCacheTime = time.Now()

	log.Printf("[Predictions] API call: %dms, %d trains\n", fetchDuration.Milliseconds(), len(resp.Trains))

	return cachedPredictions, nil
}

// startBackgroundRefresh starts a background loop to refresh data at specified intervals
func startBackgroundRefresh(name string, interval time.Duration, refreshFunc func() error) {
	// Run immediately on startup
	if err := refreshFunc(); err != nil {
		log.Printf("ERROR: %s initial refresh failed: %v\n", name, err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := refreshFunc(); err != nil {
			log.Printf("ERROR: %s refresh failed: %v\n", name, err)
		}
	}
}
