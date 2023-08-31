package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type StatusResult struct {
	URL   string `json:"url"`
	Code  int    `json:"status_code"`
	Error string `json:"error,omitempty"`
}

func main() {
	var wg sync.WaitGroup
	var mut sync.Mutex

	fileName := flag.String("f", "", "Name of the input file")
	outputFileName := flag.String("o", "", "Name of the output JSON file")
	flag.Parse()

	data, err := ioutil.ReadFile(*fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	var results []StatusResult

	for _, url := range lines {
		wg.Add(1)
		go getStatusCode(url, &wg, &mut, &results)
	}

	wg.Wait()

	if *outputFileName != "" {
		if err := saveJSONToFile(results, *outputFileName); err != nil {
			fmt.Println("Error saving JSON to file:", err)
			return
		}
	} else {
		printJSON(results)
	}
	printJSON(results)
}

func getStatusCode(url string, w *sync.WaitGroup, mut *sync.Mutex, results *[]StatusResult) {
	defer w.Done()
	resp, err := http.Get(url)

	result := StatusResult{URL: url}

	if err != nil {
		result.Error = "Error in reaching"
	} else {
		result.Code = resp.StatusCode
	}
	printJSON(results)
	mut.Lock()
	*results = append(*results, result)
	mut.Unlock()
}

func saveJSONToFile(data interface{}, fileName string) error {
	jsonBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, jsonBytes, 0644)
}

func printJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}
