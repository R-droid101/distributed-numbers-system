package main

import (
	"net/http"
	"sync"
	"testing"
)

func TestSimulatePublishers(t *testing.T) {
	var wg sync.WaitGroup
	publishers := []string{
		"http://localhost:8081/publish",
		"http://localhost:8082/publish",
		"http://localhost:8083/publish",
		"http://localhost:8084/publish",
		"http://localhost:8085/publish",
	}

	for _, url := range publishers {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			req, _ := http.NewRequest("POST", u, nil)
			req.Header.Set("Authorization", "Bearer mysecuretoken")
			http.DefaultClient.Do(req)
		}(url)
	}

	wg.Wait()
}
