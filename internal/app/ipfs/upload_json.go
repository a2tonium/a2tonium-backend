package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

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
