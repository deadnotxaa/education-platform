package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	urls := []string{
		"http://backend:8081/query1",
		"http://backend:8081/query2",
		"http://backend:8081/query3",
	}

	counter := uint64(0)

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	for {
		url := urls[counter%uint64(len(urls))]

		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", url, err)
			time.Sleep(1 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading body: %v\n", err)
		}

		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
			continue
		}

		counter++
		if counter%100 == 0 {
			fmt.Printf("Executed: %d queries\n", counter)
			fmt.Printf("Last response: %s\n", string(body))
		}

		time.Sleep(time.Duration(50+rand.Intn(50)) * time.Millisecond)
	}
}