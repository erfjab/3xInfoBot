package panel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var DefaultClient *Client

type Client struct {
	baseURL    string
	panelPath  string
	httpClient *http.Client
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

type ClientTrafficResponse struct {
	Success bool            `json:"success"`
	Msg     string          `json:"msg"`
	Obj     []ClientTraffic `json:"obj"`
}

type ClientTraffic struct {
	ID         int    `json:"id"`
	InboundID  int    `json:"inboundId"`
	Enable     bool   `json:"enable"`
	Email      string `json:"email"`
	Up         int64  `json:"up"`
	Down       int64  `json:"down"`
	ExpiryTime int64  `json:"expiryTime"`
	Total      int64  `json:"total"`
}

func NewClient(baseURL, panelPath string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}
	return &Client{
		baseURL:   baseURL,
		panelPath: panelPath,
		httpClient: &http.Client{
			Jar: jar,
		},
	}, nil
}

func (c *Client) Login(username, password, twoFactorCode string) error {
	loginURL := fmt.Sprintf("%s/%s/login/", c.baseURL, c.panelPath)

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("username", username)
	_ = writer.WriteField("password", password)
	_ = writer.WriteField("twoFactorCode", twoFactorCode)
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to build login payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, loginURL, payload)
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}

	if !loginResp.Success {
		return fmt.Errorf("login failed: %s", loginResp.Msg)
	}

	return nil
}

func (c *Client) GetClientTrafficByUUID(uuid string) (*ClientTrafficResponse, error) {
	if _, err := url.Parse(uuid); err != nil || uuid == "" {
		return nil, fmt.Errorf("invalid uuid")
	}

	apiURL := fmt.Sprintf("%s/%s/panel/api/inbounds/getClientTrafficsById/%s",
		c.baseURL, c.panelPath, url.PathEscape(uuid))

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create traffic request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("traffic request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read traffic response: %w", err)
	}

	var trafficResp ClientTrafficResponse
	if err := json.Unmarshal(body, &trafficResp); err != nil {
		return nil, fmt.Errorf("failed to parse traffic response: %w", err)
	}

	return &trafficResp, nil
}
