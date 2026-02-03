package f5xc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	BaseURL = "https://%s.console.ves.volterra.io/api"
)

func NewClient(tenantName string, apiKey string) *F5XCClient {
	return &F5XCClient{
		BaseURL: fmt.Sprintf(BaseURL, tenantName),
		ApiKey:  apiKey,
		Client:  &http.Client{},
	}
}

func (c *F5XCClient) send(req *http.Request, resData interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("APIToken %s", c.ApiKey))

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return err
	}

	return nil
}
