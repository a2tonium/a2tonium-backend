package ipfs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/a2tonium/a2tonium-backend/pkg/logger"
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
func FetchCourseMetadata(ctx context.Context, ipfsURI string) (*CourseMetadata, error) {
	if !strings.HasPrefix(ipfsURI, "ipfs://") {
		logger.ErrorKV(ctx, "ipfs.FetchQuiz failed:", logger.Err, "invalid IPFS URI")
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

	var metadata CourseMetadata
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
