package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type StatusResult struct {
	URL   string `json:"url"`
	Code  int    `json:"status_code"`
	Error string `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/check-status", handleStatusCheck)
	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handleStatusCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var urls []string
	if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	var mut sync.Mutex
	var results []StatusResult

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			result := getStatusCode(url)
			mut.Lock()
			results = append(results, result)
			mut.Unlock()
		}(url)
	}

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func getStatusCode(url string) StatusResult {
	resp, err := http.Get(url)
	result := StatusResult{URL: url}

	if err != nil {
		result.Error = "Error reaching the URL"
		return result
	}
	defer resp.Body.Close()

	result.Code = resp.StatusCode
	return result
}
