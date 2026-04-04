package upload

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotConfigured = errors.New("cloudinary is not configured")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUploadFailed  = errors.New("upload failed")
)

type UploadResult struct {
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
}

type Service interface {
	UploadImage(ctx context.Context, file multipart.File, filename string) (UploadResult, error)
}

type service struct {
	cloudName string
	apiKey    string
	apiSecret string
	folder    string
	client    *http.Client
}

func NewService(cloudName, apiKey, apiSecret, folder string) Service {
	return &service{
		cloudName: strings.TrimSpace(cloudName),
		apiKey:    strings.TrimSpace(apiKey),
		apiSecret: strings.TrimSpace(apiSecret),
		folder:    strings.Trim(strings.TrimSpace(folder), "/"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *service) UploadImage(ctx context.Context, file multipart.File, filename string) (UploadResult, error) {
	if strings.TrimSpace(filename) == "" {
		return UploadResult{}, ErrInvalidInput
	}
	if s.cloudName == "" || s.apiKey == "" || s.apiSecret == "" {
		return UploadResult{}, ErrNotConfigured
	}

	timestamp := time.Now().UTC().Unix()
	signature := s.buildSignature(timestamp)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("api_key", s.apiKey); err != nil {
		return UploadResult{}, err
	}
	if err := writer.WriteField("timestamp", strconv.FormatInt(timestamp, 10)); err != nil {
		return UploadResult{}, err
	}
	if err := writer.WriteField("signature", signature); err != nil {
		return UploadResult{}, err
	}
	if s.folder != "" {
		if err := writer.WriteField("folder", s.folder); err != nil {
			return UploadResult{}, err
		}
	}

	filePart, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return UploadResult{}, err
	}
	if _, err := io.Copy(filePart, file); err != nil {
		return UploadResult{}, err
	}
	if err := writer.Close(); err != nil {
		return UploadResult{}, err
	}

	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", s.cloudName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return UploadResult{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return UploadResult{}, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return UploadResult{}, err
	}

	var cloudResp struct {
		SecureURL string `json:"secure_url"`
		URL       string `json:"url"`
		PublicID  string `json:"public_id"`
		Error     struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &cloudResp); err != nil {
		return UploadResult{}, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if cloudResp.Error.Message != "" {
			return UploadResult{}, fmt.Errorf("%w: %s", ErrUploadFailed, cloudResp.Error.Message)
		}
		return UploadResult{}, ErrUploadFailed
	}

	url := strings.TrimSpace(cloudResp.SecureURL)
	if url == "" {
		url = strings.TrimSpace(cloudResp.URL)
	}
	if url == "" {
		return UploadResult{}, ErrUploadFailed
	}

	return UploadResult{
		URL:      url,
		PublicID: cloudResp.PublicID,
	}, nil
}

func (s *service) buildSignature(timestamp int64) string {
	params := "timestamp=" + strconv.FormatInt(timestamp, 10)
	if s.folder != "" {
		params = "folder=" + s.folder + "&" + params
	}

	hash := sha1.Sum([]byte(params + s.apiSecret))
	return hex.EncodeToString(hash[:])
}
