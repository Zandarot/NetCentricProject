package main

import (
	"fmt"
	"os"

	"MangaHub/cmd/cli/commands"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "register":
		if len(os.Args) != 4 {
			fmt.Println("Usage: mangahub register <username> <password>")
			return
		}
		commands.Register(os.Args[2], os.Args[3])
	case "login":
		if len(os.Args) != 4 {
			fmt.Println("Usage: mangahub login <username> <password>")
			return
		}
		commands.Login(os.Args[2], os.Args[3])
	case "logout":
		commands.Logout()
	case "list-manga", "manga":
		commands.ListManga()
	case "search":
		if len(os.Args) != 3 {
			fmt.Println("Usage: mangahub search <keyword>")
			return
		}
		commands.SearchManga(os.Args[2])
	case "details":
		if len(os.Args) != 3 {
			fmt.Println("Usage: mangahub details <manga_id>")
			return
		}
		commands.MangaDetails(os.Args[2])
	case "add":
		if len(os.Args) != 3 {
			fmt.Println("Usage: mangahub add <manga_id>")
			return
		}
		commands.AddToLibrary(os.Args[2])
	case "progress":
		if len(os.Args) == 2 {
			commands.GetProgress()
		} else if len(os.Args) == 4 && os.Args[2] == "update" {
			commands.UpdateProgress(os.Args[3])
		} else {
			fmt.Println("Usage: mangahub progress")
			fmt.Println("       mangahub progress update <manga_id>")
		}
	case "chat":
		commands.Chat()
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printHelp()
	}
}

func printHelp() {
	fmt.Println("MangaHub CLI Client")
	fmt.Println("====================")
	fmt.Println("Commands:")
	fmt.Println("  register <username> <password>  - Register new user")
	fmt.Println("  login <username> <password>     - Login and save token")
	fmt.Println("  logout                          - Remove saved token")
	fmt.Println("  manga | list-manga              - List all manga")
	fmt.Println("  search <keyword>                - Search manga by title/genre")
	fmt.Println("  details <manga_id>              - Get manga details")
	fmt.Println("  add <manga_id>                  - Add manga to library")
	fmt.Println("  progress                        - Show your reading progress")
	fmt.Println("  progress update <manga_id>      - Update progress for manga")
	fmt.Println("  chat                            - Join chat room")
	fmt.Println("  help                            - Show this help")
	fmt.Println("\nExamples:")
	fmt.Println("  mangahub register user1 pass123")
	fmt.Println("  mangahub login user1 pass123")
	fmt.Println("  mangahub search One")
	fmt.Println("  mangahub add m1")
	fmt.Println("  mangahub progress update m1")
}
