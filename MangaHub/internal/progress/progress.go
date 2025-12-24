package progress

import (
	"MangaHub/internal/tcp"
	"MangaHub/internal/udp"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateRequest struct {
	CurrentChapter int `json:"current_chapter"`
}

func UpdateProgress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		mangaID := c.Param("manga_id")

		var req UpdateRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
			return
		}

		// Get current total chapters to check if user completed the manga
		var totalChapters int
		db.QueryRow(`SELECT total_chapters FROM manga WHERE id = ?`, mangaID).Scan(&totalChapters)

		_, err = db.Exec(
			`UPDATE user_progress
			 SET current_chapter = ?, updated_at = CURRENT_TIMESTAMP
			 WHERE user_id = ? AND manga_id = ?`,
			req.CurrentChapter, userID, mangaID,
		)

		// Send TCP sync for user's devices
		tcp.SendProgressSync(userID, mangaID, req.CurrentChapter)

		// Broadcast via UDP if chapter is new (optional feature)
		// You could add logic here to broadcast when new chapters are released
		if req.CurrentChapter >= totalChapters && totalChapters > 0 {
			// User completed the manga - notify via UDP
			udp.BroadcastNewChapter(mangaID, req.CurrentChapter)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Cannot update progress",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Progress updated",
		})
	}
}
