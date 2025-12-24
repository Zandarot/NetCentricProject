package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Register(username, password string) {
	body := map[string]string{
		"username": username,
		"password": password,
	}

	data, _ := json.Marshal(body)

	resp, err := http.Post(apiBase+"/auth/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Registration failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 201 {
		fmt.Println("✓ Registration successful! You can now login.")
	} else {
		fmt.Printf("✗ Registration failed (%d)\n", resp.StatusCode)
		fmt.Println(string(responseBody))
	}
}
