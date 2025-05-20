package ipfs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	Count int   `json:"count"`
	Rows  []Row `json:"rows"`
}

type Row struct {
	ID            string   `json:"id"`
	CreatedAt     string   `json:"created_at"`
	Domain        string   `json:"domain"`
	Restrict      bool     `json:"restrict"`
	CustomDomains []string `json:"custom_domains"`
}

func takePinataGateway() (string, error) {
	url := "https://api.pinata.cloud/v3/ipfs/gateways"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySW5mb3JtYXRpb24iOnsiaWQiOiI4OGM1NmU4YS03OGQwLTRkYTAtYTEwMy1kZjliMmYxNjU0YTUiLCJlbWFpbCI6ImVwaWNnYW1lc3R3b0BnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicGluX3BvbGljeSI6eyJyZWdpb25zIjpbeyJkZXNpcmVkUmVwbGljYXRpb25Db3VudCI6MSwiaWQiOiJGUkExIn0seyJkZXNpcmVkUmVwbGljYXRpb25Db3VudCI6MSwiaWQiOiJOWUMxIn1dLCJ2ZXJzaW9uIjoxfSwibWZhX2VuYWJsZWQiOmZhbHNlLCJzdGF0dXMiOiJBQ1RJVkUifSwiYXV0aGVudGljYXRpb25UeXBlIjoic2NvcGVkS2V5Iiwic2NvcGVkS2V5S2V5IjoiZTEzNWUyOGIyZTA2N2RmZjQ5ODYiLCJzY29wZWRLZXlTZWNyZXQiOiI5MjZhNTc2MGUwNWEzZGFjZDZmY2FjMmZhNjk4ZTg2YzdkM2I5MTNkZGExYzdhNGNhYjYzYTIxYzk1NmE2MjE1IiwiZXhwIjoxNzc4OTYzNzgwfQ.gonxfkbUR6YqA-p93o7AKul8O9enyHf8m0h5qXq9Hsg")
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Non-OK HTTP status: %s\nResponse body: %s", resp.Status, string(bodyBytes))
	}

	// Parse the JSON response
	var response Response
	var domain string
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		log.Fatalf("Error parsing JSON response: %v", err)
	}

	// Extract the domain from the first row if available
	if len(response.Data.Rows) > 0 {
		domain = response.Data.Rows[0].Domain
		fmt.Printf("Domain: %s\n", domain)
	} else {
		fmt.Println("No rows found in response data")
	}

	return domain, nil
}
