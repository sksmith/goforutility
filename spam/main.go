package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const url = "http://localhost:8080"

// Example is a simple example response object
type Example struct {
	Duration int `json:"duration"`
}

func main() {
	count, err := strconv.Atoi(os.Args[1])
	handle(err)

	c := make(chan int, count)

	for i := 0; i < count; i++ {
		go request(c)
	}

	max := -1
	sum := -1
	for i := 0; i < count; i++ {
		v := <-c
		if max < v {
			max = v
		}
		sum += v
	}
	avg := sum / count

	log.Printf("cnt=[%d] max=[%d] sum=[%d] avg=[%d]\n", count, max, sum, avg)
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
