package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/a2tonium/a2tonium-backend/internal/infrastructure/filename"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// PinataUploadResponse models the relevant part of Pinata's response JSON
type PinataUploadResponse struct {
	Cid string `json:"cid"` // CID of uploaded file
	// other fields omitted
}

// UploadJSONToPinata uploads the JSON string as a file to Pinata and returns the IPFS CID.
func (i *IpfsService) UploadJSONToPinata(ctx context.Context, jsonData string) (string, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add file part
	part, err := writer.CreateFormFile("file", filename.GetCertificateName())
	if err != nil {
		return "", fmt.Errorf("create form file error: %w", err)
	}

	_, err = io.Copy(part, strings.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("copy json data error: %w", err)
	}

	err = writer.WriteField("network", NetworkPublic)
	if err != nil {
		return "", fmt.Errorf("write field error: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("close writer error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, PinataBaseURL, &requestBody)
	if err != nil {
		return "", fmt.Errorf("new request error: %w", err)
	}
	req.Header.Set(ContentTypeHeader, writer.FormDataContentType())
	req.Header.Set(AuthorizationHeader, "Bearer "+i.jwt)

	resp, err := i.client.Do(req)
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
	var pinataResp Response
	err = json.Unmarshal(respBody, &pinataResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return pinataResp.Data.Cid, nil
}
