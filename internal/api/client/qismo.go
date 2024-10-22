package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url" 
	"os"
	"strings"
)

const (
	AplicationJSON = "application/json"
	URLEncoded     = "application/x-www-form-urlencoded"
)

type QismoClient struct {
	BaseURL    string
	AppID      string
	SecretKey  string
	httpClient *http.Client
}

func NewQismoClient(baseURL string, appID string, secretKey string, httpClient *http.Client) *QismoClient {
	return &QismoClient{
		BaseURL:    baseURL,
		AppID:      appID,
		SecretKey:  secretKey,
		httpClient: httpClient,
	}
}

// CallAPI makes a generic API call
func (c *QismoClient) CallAPI(method, endpoint string, contentType string, body interface{}, headers map[string]string, response interface{}) error {
	// Construct the full URL
	fullURL := c.BaseURL + endpoint

	// Marshal the body if it's not nil
	var bodyReader io.Reader
	if body != nil {
		switch contentType {
		case "application/json":
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("error marshaling JSON body: %w", err)
			}
			bodyReader = bytes.NewBuffer(jsonBody)

		case "application/x-www-form-urlencoded":
			data := body.(map[string]string)
			form := url.Values{}
			for key, value := range data {
				form.Set(key, value)
			}
			bodyReader = strings.NewReader(form.Encode())

		case "multipart/form-data":
			// Assuming body is a map of form values and files
			var b bytes.Buffer
			writer := multipart.NewWriter(&b)

			for key, value := range body.(map[string]interface{}) {
				switch val := value.(type) {
				case string:
					_ = writer.WriteField(key, val)
				case *os.File:
					part, err := writer.CreateFormFile(key, val.Name())
					if err != nil {
						return fmt.Errorf("error creating form file: %w", err)
					}
					_, err = io.Copy(part, val)
					if err != nil {
						return fmt.Errorf("error writing file to form: %w", err)
					}
				}
			}

			err := writer.Close()
			if err != nil {
				return fmt.Errorf("error closing multipart writer: %w", err)
			}
			bodyReader = &b
			headers["Content-Type"] = writer.FormDataContentType()
		}
	}

	// Create the request
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Qiscus-App-Id", c.AppID)
	req.Header.Set("Qiscus-Secret-Key", c.SecretKey)
	req.Header.Set("Content-Type", contentType)

	// Set default Content-Type if not provided
	if _, ok := headers["Content-Type"]; !ok {
		req.Header.Set("Content-Type", "application/json")
	}

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-200 status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Unmarshal the response body into the provided struct
	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("error decoding response body: %w", err)
		}
	}

	return nil
}
