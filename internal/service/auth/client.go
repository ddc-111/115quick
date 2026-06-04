package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	BaseAPIURL     = "https://proapi.115.com"
	PassportAPIURL = "https://passportapi.115.com"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	ErrNo   int    `json:"errno"`
	ErrorStr string `json:"error"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("api error: code=%d message=%s", e.Code, e.Message)
	}
	return fmt.Sprintf("api error: code=%d", e.Code)
}

func (e *APIError) IsTokenError() bool {
	return e.Code == 401 || e.Code == 40002 || e.Code == 40140126
}

type Client struct {
	tokenMgr   *Manager
	httpClient *http.Client
}

func NewClient(tokenMgr *Manager) *Client {
	return &Client{
		tokenMgr: tokenMgr,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type Response struct {
	State    bool            `json:"state"`
	Code     int             `json:"code"`
	Message  string          `json:"message"`
	Data     json.RawMessage `json:"data"`
	ErrorStr string          `json:"error"`
	ErrNo    int             `json:"errno"`
}

func (c *Client) Get(path string, params map[string]string) (*Response, error) {
	return c.doRequest("GET", path, params, nil)
}

func (c *Client) Post(path string, params map[string]string, body io.Reader) (*Response, error) {
	return c.doRequest("POST", path, params, body)
}

func (c *Client) doRequest(method, path string, params map[string]string, body io.Reader) (*Response, error) {
	token, err := c.tokenMgr.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	u := BaseAPIURL + path
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		if method == "GET" {
			u += "?" + values.Encode()
		}
	}

	var req *http.Request
	if method == "GET" {
		req, err = http.NewRequest(method, u, nil)
	} else {
		if body == nil && len(params) > 0 {
			values := url.Values{}
			for k, v := range params {
				values.Set(k, v)
			}
			req, err = http.NewRequest(method, u, strings.NewReader(values.Encode()))
			if err == nil {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
		} else {
			req, err = http.NewRequest(method, u, body)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if !apiResp.State {
		apiErr := &APIError{
			Code:    apiResp.Code,
			Message: apiResp.Message,
			ErrNo:   apiResp.ErrNo,
			ErrorStr: apiResp.ErrorStr,
		}

		if apiErr.IsTokenError() {
			if refreshErr := c.tokenMgr.ForceRefresh(); refreshErr == nil {
				return c.doRequest(method, path, params, body)
			}
		}

		return nil, apiErr
	}

	return &apiResp, nil
}
