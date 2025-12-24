Run ALL the main.go in these folder api_server,grpc_server,tcp_server,ws_server to activate all server.
###FOR CLIENT TEST###

# 1. Register & login
go run cmd/cli/main.go register demo demo123
go run cmd/cli/main.go login demo demo123

# 2. Browse manga
go run cmd/cli/main.go list-manga
go run cmd/cli/main.go details m1

# 3. List all manga
go run cmd/cli/main.go list-manga

# 4. Search manga
go run cmd/cli/main.go search One

# 5. Get manga details
go run cmd/cli/main.go details m1

# 6. Add to library & update progress
go run cmd/cli/main.go add m1
go run cmd/cli/main.go progress update m1
# Enter: 100

# 7. Show progress
go run cmd/cli/main.go progress

# 8. Join chat (in another terminal)
go run cmd/cli/main.go chat
# Type messages (for private message: Type: @"USERNAME" _MESSAGES_)


