package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Helper function to set CORS headers and handle preflight requests
// CORS = Cross-Origin Resource Sharing. Browsers block requests between different origins (different ports/domains) for security.
// Our frontend runs on localhost:3000, backend on localhost:8080, different origins.
// These headers tell the browser: "Hey, it's okay, I allow localhost:3000 to access my data."
// Since we're serving frontend and API from the same origin (localhost:8080), CORS is not needed anymore
// But we keep minimal headers for compatibility.
func handleCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin (safe since we're serving everything from :8080). Was previously http://localhost:3000
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

// Helper function to write JSON responses
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Helper function to write error responses
func writeError(w http.ResponseWriter, msg string, code int) {
	http.Error(w, msg, code)
}

// Generic handler wrapper (reduces boilerplate in handlers)
func apiHandler(apiKey string, handler func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handleCORS(w, r) {
			return
		}
		handler(w, r, apiKey)
	}
}

func registerHandlers(apiKey string) {
	// Handler for /stations
	http.HandleFunc("/stations", apiHandler(apiKey, func(w http.ResponseWriter, r *http.Request, key string) {
		detailedStations, err := fetchAllStations(key)
		if err != nil {
			log.Println("ERROR /stations:", err)
			writeError(w, "API fetch failed", 500)
			return
		}
		writeJSON(w, detailedStations)
	}))

	// Handler for /entrances
	http.HandleFunc("/entrances", apiHandler(apiKey, func(w http.ResponseWriter, r *http.Request, key string) {
		// Query param: ?code=STATIONCODE. This lets the frontend request entrances for just one station,
		// so we filter the big array on the backend and only send relevant entrances.
		// This saves bandwidth and keeps the frontend simple.
		stationCode := r.URL.Query().Get("code")
		if stationCode == "" {
			writeError(w, "Missing station code", 400)
			return
		}

		// Ensure cache is populated
		if _, err := fetchAllStations(apiKey); err != nil {
			log.Println("ERROR /entrances:", err)
			writeError(w, "Cache fetch failed", 500)
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

		writeJSON(w, stationEntrances)
	}))

	// Handler for /nexttrains
	http.HandleFunc("/nexttrains", apiHandler(apiKey, func(w http.ResponseWriter, r *http.Request, key string) {
		predictions, err := fetchTrainPredictions(key)
		if err != nil {
			log.Println("ERROR /nexttrains:", err)
			writeError(w, "API fetch failed", 500)
			return
		}
		writeJSON(w, predictions)
	}))

	// Handler for /lines
	http.HandleFunc("/lines", apiHandler(apiKey, func(w http.ResponseWriter, r *http.Request, key string) {
		if _, err := fetchAllStations(apiKey); err != nil {
			log.Println("ERROR /lines:", err)
			writeError(w, "Cache fetch failed", 500)
			return
		}
		cacheMutex.RLock()
		lines := cachedLines
		cacheMutex.RUnlock()
		writeJSON(w, lines)
	}))

	// Handler for /parking
	http.HandleFunc("/parking", apiHandler(apiKey, func(w http.ResponseWriter, r *http.Request, key string) {
		stationCode := r.URL.Query().Get("code")

		if _, err := fetchAllStations(apiKey); err != nil {
			log.Println("ERROR /parking:", err)
			writeError(w, "Cache fetch failed", 500)
			return
		}

		cacheMutex.RLock()
		parking := cachedParking
		cacheMutex.RUnlock()

		// If a station code is provided, filter for that station
		if stationCode != "" {
			for _, p := range parking {
				if p.Code == stationCode {
					writeJSON(w, p)
					return
				}
			}
			// Not found
			writeError(w, "No parking info for that station", 404)
			return
		}

		// Otherwise, return all parking info
		writeJSON(w, parking)
	}))

	// Handler for /geojson/stations - serves static GeoJSON file for station info
	http.HandleFunc("/geojson/stations", apiHandler(apiKey, func(w http.ResponseWriter, r *http.Request, key string) {
		http.ServeFile(w, r, "Metro_Rail_Stations.geojson")
	}))

	// Handler for /geojson/lines - serves static GeoJSON file for rail lines
	http.HandleFunc("/geojson/lines", apiHandler(apiKey, func(w http.ResponseWriter, r *http.Request, key string) {
		http.ServeFile(w, r, "Metro_Rail_Lines.geojson")
	}))
}
