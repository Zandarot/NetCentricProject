package main

import (
	"MangaHub/internal/auth"
	grpcclient "MangaHub/internal/grpc"
	"MangaHub/internal/manga"
	"MangaHub/internal/progress"
	"MangaHub/pkg/database"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// create router
	r := gin.Default()

	//connect to the db
	db := database.InitDB("mangahub.db")
	err := database.TableCreate(db)
	if err != nil {
		log.Fatal("Unable to create database", err)
	}
	// After CreateTables
	if err := manga.LoadManga(db, "C:/Users/ITITICS/Subjects-2025-S2/NetCentric Lab/MangaHub/data/manga.json"); err != nil {
		log.Fatal(err)
	}

	//call auth
	r.POST("/auth/register", auth.Register(db))
	r.POST("/auth/login", auth.Login(db))

	//call to find manga

	r.GET("/manga", auth.AuthMiddleware(), manga.Search(db))
	r.GET("/manga/:id", auth.AuthMiddleware(), manga.Detail(db))

	//add and updtae user progress
	r.POST(
		"/library/:manga_id", auth.AuthMiddleware(), progress.AddToLibrary(db),
	)

	r.PUT(
		"/progress/:manga_id", auth.AuthMiddleware(), progress.UpdateProgress(db),
	)
	r.GET("/grpc/progress", auth.AuthMiddleware(), func(c *gin.Context) {
		userID := c.GetString("user_id")

		resp, err := grpcclient.GetUserProgress(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": "gRPC error"})
			return
		}

		c.JSON(200, resp)
	})

	//run server
	r.Run(":8080")
	//hTTP SERVER

}
