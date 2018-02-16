package main

import (
	"log"
	"net/http"
)

func requestChainFrom(url string) (*http.Response, error) {
	var req *http.Request
	var err error

	req, err = http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println("Error creating request: ", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	return client.Do(req)
}
