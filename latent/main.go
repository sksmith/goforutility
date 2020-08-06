package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Open the empty directory
// Create a go file, and go mod init github.com/sksmith/present
// Restart editor

// env GOOS=windows GOARCH=amd64 go build .
// env GOOS=linux GOARCH=amd64 go build .
// env GOOS=darwin GOARCH=amd64 go build .
// env GOOS=linux GOARCH=arm go build .

const port = ":8080"
const sleepmax = 2000

func handler(w http.ResponseWriter, r *http.Request) {
	n := rand.Intn(sleepmax)
	time.Sleep(time.Duration(n) * time.Millisecond)

	res, _ := json.Marshal(Example{n})

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// Example is a simple example response object
type Example struct {
	Duration int `json:"duration"`
}

func main() {
	http.HandleFunc("/", handler)

	rand.Seed(time.Now().UnixNano())

	log.Println("listening on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
