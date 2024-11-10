package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

// loginResponse represents the structure of the JSON response from the login endpoint
type loginResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

// login sends a POST request with username, password, and oem to the login URL and returns the token
func login(username, password, oem, loginURL string) (string, *http.Client, error) {
	// Create a cookie jar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", nil, err
	}

	// Create an HTTP client that uses the cookie jar
	client := &http.Client{Jar: jar}

	// Prepare form data with username, password, and oem properties
	formData := url.Values{}
	formData.Set("username", username)
	formData.Set("password", password)
	formData.Set("oem", oem)

	// Create a new POST request with the form data
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", nil, err
	}

	// Set the appropriate headers for the login request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the login request
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	// Check if the login was successful
	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("login failed: %s", resp.Status)
	}

	// Parse the JSON response to extract the token
	var loginResp loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		log.Println(&loginResp)

		return "", nil, fmt.Errorf("error decoding login response: %v", err)
	}

	// Return the extracted token and the client with cookies stored
	return loginResp.Data.Token, client, nil
}

// callAPI makes a request to the API URL using the token for authentication
func callAPI(client *http.Client, apiURL, token string) error {
	// Create a new request to the API
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}

	// Set the Authorization header with the Bearer token
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the API request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed: %s", resp.Status)
	}

	// Print the response status for demonstration
	fmt.Println("API response:", resp.Status)
	return nil
}

func main() {
	// Define the login URL, API URL, and user credentials
	loginURL := "https://api.hypon.cloud/v2/login"
	apiURL := "https://api.hypon.cloud/v2/administrator/adminInfo?refresh=true"

	username := os.Getenv("HYPON_USER")
	password := os.Getenv("HYPON_PASS")
	oem := ""

	// Log in and get the token and client with cookies
	token, client, err := login(username, password, oem, loginURL)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	// Call the API with the authenticated client and Bearer token
	if err := callAPI(client, apiURL, token); err != nil {
		log.Fatalf("API request failed: %v", err)
	}
}

// https://api.hypon.cloud/v2/plant/list2?page=1&page_size=10&refresh=true  //plant list / solar installs
