package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const concurrency = 10
const url = "https://blue-legacy-co-auth0-sink.apps.ap01.pcf.dcsg.com/auths/legacy"

var chdupes = make(chan string)
var dupes = make(map[string]bool)

// Auth represents a single member for which there can be many orders
type Auth struct {
	UserID     string `json:"userId"`
	AuthZeroID string `json:"authZeroId"`
	IdentityID string `json:"identityId"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("filename is required")
		return
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	processCSV(file, concurrency, processRecord)
}

func processCSV(r io.Reader, concurrency int, fn func(Auth)) {
	cr := csv.NewReader(r)
	sem := make(chan bool, concurrency)

	for {
		authdata, err := cr.Read()
		if err == io.EOF {
			break
		}
		authrec := parseRecord(authdata[1])
		sem <- true
		go func(record Auth) {
			defer func() { <-sem }()
			fn(record)
		}(*authrec)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

func processDupe(auth Auth) {
	fmt.Printf("processing %s\n", auth.UserID)
}

func processRecord(auth Auth) {
	defer timer(time.Now(), auth.UserID)

	js, err := json.Marshal(auth)
	if err != nil {
		fmt.Printf("error marshalling record userId=[%s] error=[%v]\n", auth.UserID, err)
	}

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewReader(js))
	if err != nil {
		fmt.Printf("error posting record userId=[%s] error=[%v]\n", auth.UserID, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Printf("unexpected response userId=[%s] statusCode=[%d]\n",
			auth.UserID, resp.StatusCode)
	}
}

func parseRecord(v string) *Auth {
	return &Auth{
		AuthZeroID: getField(v, "authZeroId"),
		IdentityID: getField(v, "identityId"),
		UserID:     getField(v, "userId"),
	}
}

func getField(s, n string) string {
	i := strings.Index(s, n)
	if i == -1 {
		return ""
	}
	i += len(n) + 2
	j := strings.Index(s[i:], "]") + i
	return s[i:j]
}

func timer(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%9s took %9d\n", name, elapsed/time.Millisecond)
}
