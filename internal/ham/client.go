package ham

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	appID      string
	appSecret  string
}

func NewClient(baseURL, appID, appSecret, certPath, keyPath string) (*Client, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("loaded x509 key pair failed: %w", err)
	}

	return &Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			},
		},
		baseURL:   baseURL,
		appID:     appID,
		appSecret: appSecret,
	}, nil
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}, target interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-ID", c.appID)
	req.Header.Set("X-App-Secret", c.appSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upstream returned status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
