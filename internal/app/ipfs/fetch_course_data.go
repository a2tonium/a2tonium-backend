package ipfs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type QuizAnswers struct {
	EncryptedAnswers string `json:"encrypted_answers"`
	SenderPublicKey  string `json:"sender_public_key"`
}

type CourseCompletion struct {
	GradeHighThan string `json:"gradeHighThan"`
	Certificate   string `json:"certificate"`
}

type CourseMetadata struct {
	Name             string             `json:"name"`
	QuizAnswers      QuizAnswers        `json:"quiz_answers"`
	CourseCompletion []CourseCompletion `json:"courseCompletion"`
}

// Fetches and decodes only the required fields from an IPFS JSON
func (i *IpfsService) FetchCourseMetadata(ctx context.Context, ipfsURI string) (*CourseMetadata, error) {
	ipfsURI = strings.TrimSpace(ipfsURI)
	if !strings.HasPrefix(ipfsURI, IpfsURIPrefix) {
		return nil, fmt.Errorf("invalid IPFS URI")
	}
	cid := strings.TrimPrefix(ipfsURI, IpfsURIPrefix)

	url := fmt.Sprintf(PinataGatewayFormat, i.gateway, cid)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add(AuthorizationHeader, "Bearer "+i.jwt)

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var metadata CourseMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &metadata, nil
}
