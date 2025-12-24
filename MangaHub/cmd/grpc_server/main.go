package main

import (
	"MangaHub/pkg/database"
	"context"
	"database/sql"
	"log"
	"net"

	pb "MangaHub/internal/grpc/pb"

	"google.golang.org/grpc"
)

type server struct {
	db *sql.DB
	pb.UnimplementedMangaServiceServer
}

func (s *server) GetUserProgress(
	ctx context.Context,
	req *pb.UserRequest,
) (*pb.ProgressResponse, error) {

	// Query actual database
	rows, err := s.db.Query(`
		SELECT manga_id, current_chapter, status 
		FROM user_progress 
		WHERE user_id = ?`, req.UserId)

	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.ProgressResponse{Progress: []*pb.Progress{}}, nil
	}
	defer rows.Close()

	var progressList []*pb.Progress
	for rows.Next() {
		var mangaID, status string
		var currentChapter int32

		err := rows.Scan(&mangaID, &currentChapter, &status)
		if err != nil {
			continue
		}

		progressList = append(progressList, &pb.Progress{
			MangaId:        mangaID,
			CurrentChapter: currentChapter,
			Status:         status,
		})
	}

	resp := &pb.ProgressResponse{
		Progress: progressList,
	}

	return resp, nil
}

func main() {
	// Connect to database
	db := database.InitDB("mangahub.db")
	defer db.Close()

	// Start gRPC server
	lis, err := net.Listen("tcp", ":9300")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMangaServiceServer(grpcServer, &server{db: db})

	log.Println("gRPC Server running on :9300 (connected to database)")
	grpcServer.Serve(lis)
}
