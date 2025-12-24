package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func AddToLibrary(mangaID string) {
	resp, err := makeRequest("POST", "/library/"+mangaID, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 201 {
		fmt.Println("✓ Manga added to your library!")
	} else {
		fmt.Printf("✗ Error (%d)\n", resp.StatusCode)
		fmt.Println(string(body))
	}
}

func GetProgress() {
	resp, err := makeRequest("GET", "/grpc/progress", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Printf("✗ Error (%d)\n", resp.StatusCode)
		fmt.Println(string(body))
		return
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	progress, ok := result["progress"].([]interface{})
	if !ok || len(progress) == 0 {
		fmt.Println(" No manga in your library yet.")
		fmt.Println("   Use 'mangahub add <manga_id>' to add manga.")
		return
	}

	fmt.Println(" Your Reading Progress:")
	for i, p := range progress {
		item := p.(map[string]interface{})
		fmt.Printf("  %d. %s\n", i+1, item["manga_id"])
		fmt.Printf("     Chapter: %v | Status: %s\n",
			item["current_chapter"], item["status"])
	}
}

func UpdateProgress(mangaID string) {
	var chapter int
	fmt.Print("Enter current chapter: ")
	fmt.Scanf("%d", &chapter)

	body := map[string]int{
		"current_chapter": chapter,
	}

	data, _ := json.Marshal(body)

	resp, err := makeRequest("PUT", "/progress/"+mangaID, data)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		fmt.Println("✓ Progress updated successfully!")
	} else {
		fmt.Printf("✗ Error (%d)\n", resp.StatusCode)
		fmt.Println(string(responseBody))
	}
}
