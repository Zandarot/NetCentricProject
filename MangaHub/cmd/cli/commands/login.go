package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Login(username, password string) {
	body := map[string]string{
		"username": username,
		"password": password,
	}

	data, _ := json.Marshal(body)

	resp, err := http.Post(apiBase+"/auth/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Printf("✗ Login failed (%d)\n", resp.StatusCode)
		fmt.Println(string(responseBody))
		return
	}

	var result map[string]interface{}
	json.Unmarshal(responseBody, &result)

	token, ok := result["TOKEN"].(string)
	if !ok || token == "" {
		fmt.Println("✗ No token in response")
		return
	}

	if err := saveToken(token); err != nil {
		fmt.Printf("✗ Failed to save token: %v\n", err)
		return
	}

	fmt.Println("✓ Login successful! Token saved.")
}
