package client

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("https://api.example.com")
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "https://api.example.com", client.BaseURL)
	assert.NotNil(t, client.HTTPClient)
	assert.Empty(t, client.Token)
}

func TestClient_Login_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v2/login", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		// Check form data
		err := r.ParseForm()
		require.NoError(t, err)
		assert.Equal(t, "testuser", r.FormValue("username"))
		assert.Equal(t, "testpass", r.FormValue("password"))
		assert.Equal(t, "testoem", r.FormValue("oem"))

		// Return successful login response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"token": "test-token-123"}}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	require.NoError(t, err)

	err = client.Login("testuser", "testpass", "testoem")
	require.NoError(t, err)
	assert.Equal(t, "test-token-123", client.Token)
}

func TestClient_Login_Failure(t *testing.T) {
	// Create a mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	require.NoError(t, err)

	err = client.Login("wronguser", "wrongpass", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "login failed")
	assert.Empty(t, client.Token)
}

func TestClient_Get_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/login" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": {"token": "test-token-123"}}`))
			return
		}

		if r.URL.Path == "/v2/test" {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "Bearer test-token-123", r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": true}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	require.NoError(t, err)

	// Login first
	err = client.Login("testuser", "testpass", "")
	require.NoError(t, err)

	// Make API call
	resp, err := client.Get("/v2/test")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestClient_Get_NotAuthenticated(t *testing.T) {
	client, err := NewClient("https://api.example.com")
	require.NoError(t, err)

	_, err = client.Get("/v2/test")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client not authenticated")
}

func TestClient_GetAdminInfo(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/login" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": {"token": "test-token-123"}}`))
			return
		}

		if strings.Contains(r.URL.Path, "/v2/administrator/adminInfo") {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "Bearer test-token-123", r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"admin": "info"}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	require.NoError(t, err)

	// Login first
	err = client.Login("testuser", "testpass", "")
	require.NoError(t, err)

	// Get admin info
	resp, err := client.GetAdminInfo()
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestClient_GetPlantList(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/login" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": {"token": "test-token-123"}}`))
			return
		}

		if strings.Contains(r.URL.Path, "/v2/plant/list2") {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "Bearer test-token-123", r.Header.Get("Authorization"))

			// Check query parameters
			query := r.URL.Query()
			assert.Equal(t, "1", query.Get("page"))
			assert.Equal(t, "10", query.Get("page_size"))
			assert.Equal(t, "true", query.Get("refresh"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"plants": []}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	require.NoError(t, err)

	// Login first
	err = client.Login("testuser", "testpass", "")
	require.NoError(t, err)

	// Get plant list
	resp, err := client.GetPlantList(1, 10)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
