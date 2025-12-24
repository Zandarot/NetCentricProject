package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Progress() {
	resp, err := http.Get("http://localhost:8080/grpc/progress")
	if err != nil {
		fmt.Println("Error")
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	fmt.Println("Your progress:")
	fmt.Println(data)
}
