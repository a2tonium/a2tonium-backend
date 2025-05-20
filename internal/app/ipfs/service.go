package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

// Structs for the required fields
type QuizAnswers struct {
	EncryptedAnswers string `json:"encrypted_answers"`
	SenderPublicKey  string `json:"sender_public_key"`
}

type CourseCompletion struct {
	GradeHighThan string `json:"gradeHighThan"`
	Certificate   string `json:"certificate"`
}

type NFTMetadata struct {
	Name             string             `json:"name"`
	QuizAnswers      QuizAnswers        `json:"quiz_answers"`
	CourseCompletion []CourseCompletion `json:"courseCompletion"`
}

func main() {

	url := "https://api.pinata.cloud/v3/files/public/%s"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer <token>")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}

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

// Fetches and decodes only the required fields from an IPFS JSON
func FetchQuizAndCompletion(ipfsURI string) (*NFTMetadata, error) {
	if !strings.HasPrefix(ipfsURI, "ipfs://") {
		return nil, fmt.Errorf("invalid IPFS URI")
	}
	cid := strings.TrimPrefix(ipfsURI, "ipfs://")

	domain, err := takePinataGateway()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s.mypinata.cloud/ipfs/%s", domain, cid)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySW5mb3JtYXRpb24iOnsiaWQiOiI4OGM1NmU4YS03OGQwLTRkYTAtYTEwMy1kZjliMmYxNjU0YTUiLCJlbWFpbCI6ImVwaWNnYW1lc3R3b0BnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicGluX3BvbGljeSI6eyJyZWdpb25zIjpbeyJkZXNpcmVkUmVwbGljYXRpb25Db3VudCI6MSwiaWQiOiJGUkExIn0seyJkZXNpcmVkUmVwbGljYXRpb25Db3VudCI6MSwiaWQiOiJOWUMxIn1dLCJ2ZXJzaW9uIjoxfSwibWZhX2VuYWJsZWQiOmZhbHNlLCJzdGF0dXMiOiJBQ1RJVkUifSwiYXV0aGVudGljYXRpb25UeXBlIjoic2NvcGVkS2V5Iiwic2NvcGVkS2V5S2V5IjoiZTEzNWUyOGIyZTA2N2RmZjQ5ODYiLCJzY29wZWRLZXlTZWNyZXQiOiI5MjZhNTc2MGUwNWEzZGFjZDZmY2FjMmZhNjk4ZTg2YzdkM2I5MTNkZGExYzdhNGNhYjYzYTIxYzk1NmE2MjE1IiwiZXhwIjoxNzc4OTYzNzgwfQ.gonxfkbUR6YqA-p93o7AKul8O9enyHf8m0h5qXq9Hsg")

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var metadata NFTMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &metadata, nil
}

// PinataUploadResponse models the relevant part of Pinata's response JSON
type PinataUploadResponse struct {
	Cid string `json:"cid"` // CID of uploaded file
	// other fields omitted
}

// UploadJSONToPinata uploads the JSON string as a file to Pinata and returns the IPFS CID.
func UploadJSONToPinata(jsonData, filename string) (string, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add file part
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("create form file error: %w", err)
	}
	_, err = io.Copy(part, bytes.NewReader([]byte(jsonData)))
	if err != nil {
		return "", fmt.Errorf("copy json data error: %w", err)
	}

	// Add form field "network"
	err = writer.WriteField("network", "public")
	if err != nil {
		return "", fmt.Errorf("write field error: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("close writer error: %w", err)
	}

	req, err := http.NewRequest("POST", "https://uploads.pinata.cloud/v3/files", &requestBody)
	if err != nil {
		return "", fmt.Errorf("new request error: %w", err)
	}
	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySW5mb3JtYXRpb24iOnsiaWQiOiI4OGM1NmU4YS03OGQwLTRkYTAtYTEwMy1kZjliMmYxNjU0YTUiLCJlbWFpbCI6ImVwaWNnYW1lc3R3b0BnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicGluX3BvbGljeSI6eyJyZWdpb25zIjpbeyJkZXNpcmVkUmVwbGljYXRpb25Db3VudCI6MSwiaWQiOiJGUkExIn0seyJkZXNpcmVkUmVwbGljYXRpb25Db3VudCI6MSwiaWQiOiJOWUMxIn1dLCJ2ZXJzaW9uIjoxfSwibWZhX2VuYWJsZWQiOmZhbHNlLCJzdGF0dXMiOiJBQ1RJVkUifSwiYXV0aGVudGljYXRpb25UeXBlIjoic2NvcGVkS2V5Iiwic2NvcGVkS2V5S2V5IjoiZTEzNWUyOGIyZTA2N2RmZjQ5ODYiLCJzY29wZWRLZXlTZWNyZXQiOiI5MjZhNTc2MGUwNWEzZGFjZDZmY2FjMmZhNjk4ZTg2YzdkM2I5MTNkZGExYzdhNGNhYjYzYTIxYzk1NmE2MjE1IiwiZXhwIjoxNzc4OTYzNzgwfQ.gonxfkbUR6YqA-p93o7AKul8O9enyHf8m0h5qXq9Hsg"
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response error: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("upload failed: status %d, response: %s", resp.StatusCode, string(respBody))
	}

	// Parse response JSON to get CID
	var pinataResp PinataUploadResponse
	err = json.Unmarshal(respBody, &pinataResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return pinataResp.Cid, nil
}
