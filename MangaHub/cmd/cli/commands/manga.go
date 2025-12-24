package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ListManga() {
	resp, err := makeRequest("GET", "/manga", nil)
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

	var mangas []map[string]interface{}
	json.Unmarshal(body, &mangas)

	fmt.Printf(" Available Manga (%d found):\n", len(mangas))
	for _, m := range mangas {
		fmt.Printf("  • %s (ID: %s)\n", m["title"], m["id"])
		fmt.Printf("    Author: %s | Status: %s | Chapters: %v\n",
			m["author"], m["status"], m["total_chapters"])
		fmt.Printf("    Genres: %s\n\n", m["genres"])
	}
}

func SearchManga(keyword string) {
	resp, err := makeRequest("GET", "/manga?search="+keyword, nil)
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

	var mangas []map[string]interface{}
	json.Unmarshal(body, &mangas)

	fmt.Printf(" Search Results for '%s' (%d found):\n", keyword, len(mangas))
	for _, m := range mangas {
		fmt.Printf("  • %s (ID: %s) - %s\n", m["title"], m["id"], m["author"])
	}
}

func MangaDetails(mangaID string) {
	resp, err := makeRequest("GET", "/manga/"+mangaID, nil)
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

	var manga map[string]interface{}
	json.Unmarshal(body, &manga)

	fmt.Println(" Manga Details:")
	fmt.Printf("  Title: %s\n", manga["title"])
	fmt.Printf("  ID: %s\n", manga["id"])
	fmt.Printf("  Author: %s\n", manga["author"])
	fmt.Printf("  Genres: %s\n", manga["genres"])
	fmt.Printf("  Status: %s\n", manga["status"])
	fmt.Printf("  Total Chapters: %v\n", manga["total_chapters"])
	fmt.Printf("  Description: %s\n", manga["description"])
}
