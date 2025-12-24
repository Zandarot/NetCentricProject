package commands

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const tokenFile = ".mangahub_token"
const apiBase = "http://localhost:8080"

// Save token to file
func saveToken(token string) error {
	return os.WriteFile(tokenFile, []byte(token), 0600)
}

// Load token from file
func loadToken() (string, error) {
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// Remove token file (logout)
func Logout() {
	os.Remove(tokenFile)
	fmt.Println("Logged out successfully.")
}

// Make authenticated request
func makeRequest(method, endpoint string, body []byte) (*http.Response, error) {
	token, err := loadToken()
	if err != nil {
		return nil, fmt.Errorf("not logged in. Please login first")
	}

	client := &http.Client{}
	url := apiBase + endpoint

	var req *http.Request
	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	return client.Do(req)
}

// ===== ADD THIS NEW FUNCTION =====
// Get username from JWT token
func getUsernameFromToken() (string, error) {
	token, err := loadToken()
	if err != nil {
		return "", err
	}

	// Simple JWT parsing (without full validation)
	// JWT format: header.payload.signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	// Decode payload (middle part)
	// Note: This is a simple decode without signature verification
	// For production, use proper JWT library
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode token: %v", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("failed to parse token claims: %v", err)
	}

	// Extract username from claims
	if username, ok := claims["username"].(string); ok {
		return username, nil
	}

	// Try alternative field names
	if username, ok := claims["user"].(string); ok {
		return username, nil
	}

	return "", fmt.Errorf("username not found in token")
}
