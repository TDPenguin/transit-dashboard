package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	// Declares err AND checks it on one line. godotenv.Load() only returns error or nil if success.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("WMATA_API_KEY")
	fmt.Println("==== Server running on :8080 ====")
	fmt.Println("Frontend: http://localhost:8080")
	fmt.Println("API: http://localhost:8080/stations")

	// Pre-warm caches sequentially on startup to avoid rate limiting
	log.Println("Pre-warming caches...")
	if _, err := refreshAllStations(apiKey); err != nil {
		log.Println("ERROR: Failed to pre-warm static cache:", err)
	}
	if _, err := refreshTrainPredictions(apiKey); err != nil {
		log.Println("ERROR: Failed to pre-warm predictions cache:", err)
	}
	log.Println("Caches pre-warmed successfully!")

	// Start background refresh loops (now that initial data is loaded)
	go startBackgroundRefresh("Predictions", 20*time.Second, func() error {
		_, err := refreshTrainPredictions(apiKey)
		return err
	})
	go startBackgroundRefresh("Static Data", 24*time.Hour, func() error {
		_, err := refreshAllStations(apiKey)
		return err
	})

	// Register API handlers
	registerHandlers(apiKey)

	// Serve frontend static files from ../frontend directory
	// This allows Go to serve index.html, script.js, style.css, etc.
	// Files are served at the root path ("/"), API handlers take precedence
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
