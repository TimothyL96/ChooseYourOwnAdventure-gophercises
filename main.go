package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type JsonStruct struct {
	Title   string
	Story   []string
	Options []struct {
		Text string
		Arc  string
	}
}

const (
	JSONFileName   string = "gopher.json"
	HTMLFileName   string = "index.html"
	JSONDefaultKey string = "intro"
)

func main() {
	// Create map to store JSON data
	jsonData := make(map[string]JsonStruct)

	// Load JSON from file
	err := LoadJSONFromFile(jsonData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// HTTP server
	http.HandleFunc("/", handleLoad(jsonData))

	log.Println("HTTP server started on port :8080")
	_ = http.ListenAndServe(":8080", nil)
}

func LoadJSONFromFile(jsonData map[string]JsonStruct) error {
	// Open file
	jsonFile, err := os.Open(JSONFileName)
	if err != nil {
		return errors.New("error: file not found: " + JSONFileName)
	}

	// Read JSON data
	jsDecoder := json.NewDecoder(jsonFile)

	for jsDecoder.More() {
		err = jsDecoder.Decode(&jsonData)

		if err != nil {
			return errors.New("error: Can't decode json: " + err.Error())
		}
	}

	return nil
}

func handleLoad(jsonData map[string]JsonStruct) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the template HTML file
		t, err := template.ParseFiles(HTMLFileName)
		if err != nil {
			fmt.Println(errors.New("error: template: " + err.Error()))
			return
		}

		// If URL is not a story arc, navigate to intro story
		jsonKey := JSONDefaultKey
		path := strings.TrimSpace(r.URL.Path[1:])
		if _, ok := jsonData[path]; ok {
			jsonKey = path
		}

		err = t.Execute(w, jsonData[jsonKey])
		if err != nil {
			errStr := "template execute error"
			log.Println(errStr, err)
			http.Error(w, errStr, http.StatusInternalServerError)
		}
	}
}
