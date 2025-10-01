package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// Client represents a Hypon API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// loginResponse represents the structure of the JSON response from the login endpoint
type loginResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

// NewClient creates a new Hypon API client
func NewClient(baseURL string) (*Client, error) {
	// Create a cookie jar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// Create an HTTP client that uses the cookie jar
	client := &http.Client{Jar: jar}

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: client,
	}, nil
}

// Login authenticates with the Hypon API and stores the token
func (c *Client) Login(username, password, oem string) error {
	loginURL := c.BaseURL + "/v2/login"

	// Prepare form data with username, password, and oem properties
	formData := url.Values{}
	formData.Set("username", username)
	formData.Set("password", password)
	formData.Set("oem", oem)

	// Create a new POST request with the form data
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	// Set the appropriate headers for the login request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the login request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the login was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	// Parse the JSON response to extract the token
	var loginResp loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("error decoding login response: %v", err)
	}

	// Store the token in the client
	c.Token = loginResp.Data.Token
	return nil
}

// Get makes a GET request to the specified API endpoint
func (c *Client) Get(endpoint string) (*http.Response, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("client not authenticated - call Login() first")
	}

	// Create a new request to the API
	req, err := http.NewRequest("GET", c.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the Bearer token
	req.Header.Set("Authorization", "Bearer "+c.Token)

	// Send the API request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed: %s", resp.Status)
	}

	return resp, nil
}

// GetAdminInfo retrieves administrator information
func (c *Client) GetAdminInfo() (*http.Response, error) {
	return c.Get("/v2/administrator/adminInfo?refresh=true")
}

// GetPlantList retrieves the list of plants/solar installations
func (c *Client) GetPlantList(page, pageSize int) (*http.Response, error) {
	endpoint := fmt.Sprintf("/v2/plant/list2?page=%d&page_size=%d&refresh=true", page, pageSize)
	return c.Get(endpoint)
}

// GetPlantListMenu retrieves the plant list from the menu action endpoint
func (c *Client) GetPlantListMenu() (*http.Response, error) {
	return c.Get("/v2/auth/menu/actionList?menu_code=PlantList&refresh=true")
}

// SetToken sets the authentication token for the client
func (c *Client) SetToken(token string) {
	c.Token = token
}
