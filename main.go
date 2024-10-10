package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// In-memory store for minimal and maximal power values
var store = struct {
	minimalPower float64
	maximalPower float64
	sync.RWMutex
}{
	minimalPower: -15,
	maximalPower: 30,
}

var body struct {
	Value float64 `json:"Value"`
}

// Handler to set minimalPower
func setMinimalPowerHandler(w http.ResponseWriter, r *http.Request) {

	// Expecting a value in the query parameters
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Printf("Invalid request method: %v\n", r.Method)
		return
	}

	var body struct {
		Value float64 `json:"value"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		log.Printf("Invalid body: %+v\n", err)
		return
	}

	// Set the minimal power safely with a write lock
	store.Lock()
	store.minimalPower = body.Value
	store.Unlock()

	fmt.Printf("Minimal Power set to: %v kW\n", body.Value)
}

// Handler to set maximalPower
func setMaximalPowerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		fmt.Printf("Invalid request method + %+v\n", r.Method)
		return
	}

	// Read and log the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		fmt.Printf("Error reading body: %v\n", err)
		return
	}

	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		fmt.Printf("Invalid body: %+v\n", err)
		return
	}

	log.Printf("Maximal Power set to: %v kvar\n", body.Value)

	// Set the maximal power safely with a write lock
	store.Lock()
	store.maximalPower = body.Value
	store.Unlock()
}

// Handler to get maximalPower
func getMaximalPowerHandler(w http.ResponseWriter, r *http.Request) {
	store.RLock()
	defer store.RUnlock()

	fmt.Printf("%v kvar\n", store.maximalPower)

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		fmt.Printf("Invalid request method: %v\n", r.Method)
		return
	}

	// Create a response with the maximal power value
	response := map[string]float64{
		"value": store.maximalPower,
	}

	// Set response header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Send JSON response
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		fmt.Printf("Error encoding response: %v\n", err)
		return
	}

}

// Handler to get minimalPower
func getMinimalPowerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		fmt.Printf("Invalid request method: %v\n", r.Method)
		return
	}

	store.RLock()
	defer store.RUnlock()

	fmt.Printf("%v kW\n", store.minimalPower)
	// Create a response with the minimal power value
	response := map[string]float64{
		"value": store.minimalPower,
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		fmt.Printf("Error encoding response: %v\n", err)
		return
	}

}

func main() {
	http.HandleFunc("/set/peakshaving/minimalPower", enableCors(setMinimalPowerHandler))
	http.HandleFunc("/get/peakshaving/minimalPower", enableCors(getMinimalPowerHandler))
	http.HandleFunc("/set/peakshaving/maximalPower", enableCors(setMaximalPowerHandler))
	http.HandleFunc("/get/peakshaving/maximalPower", enableCors(getMaximalPowerHandler))

	fmt.Println("Server is listening on port 6080...")

	log.Fatal(http.ListenAndServe(":6080", nil))
}

// Cors middleware
func enableCors(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CORS middleware")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		// (*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Handle CORS preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
