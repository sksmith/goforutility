package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const url = "http://localhost:8080"

// Example is a simple example response object
type Example struct {
	Duration int `json:"duration"`
}

func main() {
	c := make(chan int)
	go request(c)

	log.Printf("doing something while our request processes...")

	log.Printf("got: %d\n", <-c)
}

func request(c chan int) {
	response, err := http.Get(url)
	handle(err)

	body, err := ioutil.ReadAll(response.Body)
	handle(err)

	example := &Example{}
	err = json.Unmarshal(body, example)
	handle(err)

	c <- example.Duration
}

func handle(err error) {
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}
