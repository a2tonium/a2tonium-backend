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

func (i *IpfsService) getPinataGateway() (string, error) {
	req, err := http.NewRequest(http.MethodGet, PinataGatewayURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+i.jwt)

	resp, err := i.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
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
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if len(response.Data.Rows) == 0 {
		return "", fmt.Errorf("no gateway found")
	}

	return response.Data.Rows[0].Domain, nil
}
