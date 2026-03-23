package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type UtilsClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUtilsClient(baseURL string) *UtilsClient {
	return &UtilsClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *UtilsClient) UploadImage(ctx context.Context, file io.Reader, filename string, cookie string) (string, error) {
	// build multipart request
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(part, file); err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/api/utils/upload", c.baseURL), &buf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", cookie) // forward session_token cookie for auth

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("utils service returned %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			ImageURL string `json:"image_url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.ImageURL, nil
}
