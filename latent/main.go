package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const port = ":8080"
const sleepmax = 2000

func main() {
	http.HandleFunc("/", handler)

	rand.Seed(time.Now().UnixNano())

	log.Println("listening on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Example is a simple example response object
type Example struct {
	Duration int `json:"duration"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	n := rand.Intn(sleepmax)
	time.Sleep(time.Duration(n) * time.Millisecond)

	res, _ := json.Marshal(Example{n})

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
